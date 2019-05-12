package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/hive/pkg/apis/hive/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

func loadResources(path string) ([]runtime.RawExtension, error) {
	var resources = []runtime.RawExtension{}
	err := filepath.Walk(path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(p, ".yaml") {
				fmt.Println(p)
				data, err := ioutil.ReadFile(p)
				if err != nil {
					return err
				}
				jsonBytes, err := yaml.YAMLToJSON(data)
				if err != nil {
					return err
				}
				var r = runtime.RawExtension{}
				json.Unmarshal(jsonBytes, &r)
				resources = append(resources, r)
			}
			return nil
		})
	return resources, err
}

func CreateSelectorSyncSet(name string, selector string, path string) v1alpha1.SelectorSyncSet {
	resources, err := loadResources(path)
	if err != nil {
		log.Println(err)
	}

	labelSelector, err := metav1.ParseToLabelSelector(selector)
	if err != nil {
		log.Println(err)
	}

	var syncSet = &v1alpha1.SelectorSyncSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SelectorSyncSet",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"atlas.worldpay.com/generated": "true",
			},
		},
		Spec: v1alpha1.SelectorSyncSetSpec{
			SyncSetCommonSpec: v1alpha1.SyncSetCommonSpec{
				Resources:         resources,
				ResourceApplyMode: "Sync",
			},
			ClusterDeploymentSelector: *labelSelector,
		},
	}
	return *syncSet
}
func CreateSyncSet(name string, clusterName string, path string) v1alpha1.SyncSet {
	resources, err := loadResources(path)
	if err != nil {
		log.Println(err)
	}

	var syncSet = &v1alpha1.SyncSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SyncSet",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"atlas.worldpay.com/generated": "true",
			},
		},
		Spec: v1alpha1.SyncSetSpec{
			SyncSetCommonSpec: v1alpha1.SyncSetCommonSpec{
				Resources:         resources,
				ResourceApplyMode: "Sync",
			},
			ClusterDeploymentRefs: []corev1.LocalObjectReference{
				corev1.LocalObjectReference{
					Name: clusterName,
				},
			},
		},
	}
	return *syncSet
}