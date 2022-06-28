/*
Copyright 2022 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	v1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type CephClusterHandler func(string, *v1.CephCluster) (*v1.CephCluster, error)

type CephClusterController interface {
	generic.ControllerMeta
	CephClusterClient

	OnChange(ctx context.Context, name string, sync CephClusterHandler)
	OnRemove(ctx context.Context, name string, sync CephClusterHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() CephClusterCache
}

type CephClusterClient interface {
	Create(*v1.CephCluster) (*v1.CephCluster, error)
	Update(*v1.CephCluster) (*v1.CephCluster, error)
	UpdateStatus(*v1.CephCluster) (*v1.CephCluster, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.CephCluster, error)
	List(namespace string, opts metav1.ListOptions) (*v1.CephClusterList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CephCluster, err error)
}

type CephClusterCache interface {
	Get(namespace, name string) (*v1.CephCluster, error)
	List(namespace string, selector labels.Selector) ([]*v1.CephCluster, error)

	AddIndexer(indexName string, indexer CephClusterIndexer)
	GetByIndex(indexName, key string) ([]*v1.CephCluster, error)
}

type CephClusterIndexer func(obj *v1.CephCluster) ([]string, error)

type cephClusterController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewCephClusterController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) CephClusterController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &cephClusterController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromCephClusterHandlerToHandler(sync CephClusterHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.CephCluster
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.CephCluster))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *cephClusterController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.CephCluster))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateCephClusterDeepCopyOnChange(client CephClusterClient, obj *v1.CephCluster, handler func(obj *v1.CephCluster) (*v1.CephCluster, error)) (*v1.CephCluster, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *cephClusterController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *cephClusterController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *cephClusterController) OnChange(ctx context.Context, name string, sync CephClusterHandler) {
	c.AddGenericHandler(ctx, name, FromCephClusterHandlerToHandler(sync))
}

func (c *cephClusterController) OnRemove(ctx context.Context, name string, sync CephClusterHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromCephClusterHandlerToHandler(sync)))
}

func (c *cephClusterController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *cephClusterController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *cephClusterController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *cephClusterController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *cephClusterController) Cache() CephClusterCache {
	return &cephClusterCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *cephClusterController) Create(obj *v1.CephCluster) (*v1.CephCluster, error) {
	result := &v1.CephCluster{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *cephClusterController) Update(obj *v1.CephCluster) (*v1.CephCluster, error) {
	result := &v1.CephCluster{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *cephClusterController) UpdateStatus(obj *v1.CephCluster) (*v1.CephCluster, error) {
	result := &v1.CephCluster{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *cephClusterController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *cephClusterController) Get(namespace, name string, options metav1.GetOptions) (*v1.CephCluster, error) {
	result := &v1.CephCluster{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *cephClusterController) List(namespace string, opts metav1.ListOptions) (*v1.CephClusterList, error) {
	result := &v1.CephClusterList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *cephClusterController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *cephClusterController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.CephCluster, error) {
	result := &v1.CephCluster{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type cephClusterCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *cephClusterCache) Get(namespace, name string) (*v1.CephCluster, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.CephCluster), nil
}

func (c *cephClusterCache) List(namespace string, selector labels.Selector) (ret []*v1.CephCluster, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CephCluster))
	})

	return ret, err
}

func (c *cephClusterCache) AddIndexer(indexName string, indexer CephClusterIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.CephCluster))
		},
	}))
}

func (c *cephClusterCache) GetByIndex(indexName, key string) (result []*v1.CephCluster, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.CephCluster, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.CephCluster))
	}
	return result, nil
}

type CephClusterStatusHandler func(obj *v1.CephCluster, status v1.ClusterStatus) (v1.ClusterStatus, error)

type CephClusterGeneratingHandler func(obj *v1.CephCluster, status v1.ClusterStatus) ([]runtime.Object, v1.ClusterStatus, error)

func RegisterCephClusterStatusHandler(ctx context.Context, controller CephClusterController, condition condition.Cond, name string, handler CephClusterStatusHandler) {
	statusHandler := &cephClusterStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromCephClusterHandlerToHandler(statusHandler.sync))
}

func RegisterCephClusterGeneratingHandler(ctx context.Context, controller CephClusterController, apply apply.Apply,
	condition condition.Cond, name string, handler CephClusterGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &cephClusterGeneratingHandler{
		CephClusterGeneratingHandler: handler,
		apply:                        apply,
		name:                         name,
		gvk:                          controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterCephClusterStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type cephClusterStatusHandler struct {
	client    CephClusterClient
	condition condition.Cond
	handler   CephClusterStatusHandler
}

func (a *cephClusterStatusHandler) sync(key string, obj *v1.CephCluster) (*v1.CephCluster, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type cephClusterGeneratingHandler struct {
	CephClusterGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *cephClusterGeneratingHandler) Remove(key string, obj *v1.CephCluster) (*v1.CephCluster, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.CephCluster{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *cephClusterGeneratingHandler) Handle(obj *v1.CephCluster, status v1.ClusterStatus) (v1.ClusterStatus, error) {
	objs, newStatus, err := a.CephClusterGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}