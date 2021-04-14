/*
Copyright 2012 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package singleflight provides a duplicate function call suppression
// mechanism.
package singleflight

import "sync"

// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
//同一个对象多次同时多次调用这个逻辑的时候，可以使用其中的一个去执行
func (g *Group) copyDo(key string, fn func()(interface{},error)) (interface{}, error ){
	g.mu.Lock() //加锁保护存放key的map，因为要并发执行
	if g.m == nil { //lazing make 方式建立
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok { //如果map中已经存在对这个key的处理那就等着吧
	    g.mu.Unlock() //解锁，对map的操作已经完毕
		c.wg.Wait()
		return c.val,c.err //map中只有一份key，所以只有一个c
	}
	c := new(call) //创建一个工作单元，只负责处理一种key
	c.wg.Add(1)
	g.m[key] = c //将key注册到map中
	g.mu.Unlock() //map的操做完成，解锁

	c.val, c.err = fn()//第一个注册者去执行
	c.wg.Done()

	g.mu.Lock()
	delete(g.m,key) //对map进行操作，需要枷锁
	g.mu.Unlock()

	return c.val, c.err //给第一个注册者返回结果
}
