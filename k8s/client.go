/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package k8s

import (
	"edp-admin-console/service/logger"
	edppipelinesv1alpha1 "github.com/epmd-edp/cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	edpv1alpha1 "github.com/epmd-edp/codebase-operator/v2/pkg/apis/edp/v1alpha1"
	appsV1Client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	storageV1Client "k8s.io/client-go/kubernetes/typed/storage/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var log = logger.GetLogger()
var k8sConfig clientcmd.ClientConfig
var SchemeGroupVersion = schema.GroupVersion{Group: "v2.edp.epam.com", Version: "v1alpha1"}

type ClientSet struct {
	CoreClient     *coreV1Client.CoreV1Client
	StorageClient  *storageV1Client.StorageV1Client
	EDPRestClient  *rest.RESTClient
	AppsV1Client   *appsV1Client.AppsV1Client
	K8sAppV1Client v1.AppsV1Interface
}

func init() {
	k8sConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
}

func CreateOpenShiftClients() ClientSet {
	coreClient, err := getCoreClient()
	if err != nil {
		log.Error("An error has occurred while getting core client", zap.Error(err))
		panic(err)
	}

	crClient, err := getApplicationClient()
	if err != nil {
		log.Error("An error has occurred while getting custom resource client", zap.Error(err))
		panic(err)
	}

	storageClient, err := getStorageClient()
	if err != nil {
		log.Error("An error has occurred while getting custom resource client", zap.Error(err))
		panic(err)
	}

	openshiftAppClient, err := getOpenshiftApplicationClient()
	if err != nil {
		log.Error("An error has occurred while getting oenshift application resource client", zap.Error(err))
		panic(err)
	}

	k8sAppClient, err := getK8sAppsV1Client()
	if err != nil {
		log.Error("An error has occurred while getting k8s extension client", zap.Error(err))
		panic(err)
	}

	return ClientSet{
		CoreClient:     coreClient,
		StorageClient:  storageClient,
		EDPRestClient:  crClient,
		AppsV1Client:   openshiftAppClient,
		K8sAppV1Client: k8sAppClient,
	}
}

func getCoreClient() (*coreV1Client.CoreV1Client, error) {
	restConfig, err := k8sConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	coreClient, err := coreV1Client.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return coreClient, nil
}

func getK8sAppsV1Client() (v1.AppsV1Interface, error) {
	restConfig, err := k8sConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return cs.AppsV1(), nil
}

func getStorageClient() (*storageV1Client.StorageV1Client, error) {
	restConfig, err := k8sConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	coreClient, err := storageV1Client.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return coreClient, nil
}

func getApplicationClient() (*rest.RESTClient, error) {
	var config *rest.Config
	var err error

	config, err = k8sConfig.ClientConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := createCrdClient(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getOpenshiftApplicationClient() (*appsV1Client.AppsV1Client, error) {
	var config *rest.Config
	var err error

	config, err = k8sConfig.ClientConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := appsV1Client.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func createCrdClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&edpv1alpha1.Codebase{},
		&edpv1alpha1.CodebaseList{},
		&edpv1alpha1.CodebaseBranch{},
		&edpv1alpha1.CodebaseBranchList{},
		&edppipelinesv1alpha1.CDPipeline{},
		&edppipelinesv1alpha1.CDPipelineList{},
		&edppipelinesv1alpha1.Stage{},
		&edppipelinesv1alpha1.StageList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
