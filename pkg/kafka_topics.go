package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/strimzioperator/kafka/v1beta2"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/helm/convertmaps"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func kafkaTopics(ctx *pulumi.Context, locals *Locals, createdNamespace *kubernetescorev1.Namespace,
	createdKafkaCluster *v1beta2.Kafka, labels map[string]string) error {

	for _, kafkaTopic := range locals.KafkaKubernetes.Spec.KafkaTopics {

		config := vars.KafkaTopicDefaultConfig
		for k, v := range kafkaTopic.Config {
			config[k] = v
		}

		_, err := v1beta2.NewKafkaTopic(ctx,
			kafkaTopic.Name,
			&v1beta2.KafkaTopicArgs{
				Metadata: metav1.ObjectMetaArgs{
					Name:      pulumi.String(kafkaTopic.Name),
					Namespace: createdNamespace.Metadata.Name(),
					Labels:    pulumi.ToStringMap(labels),
				},
				Spec: v1beta2.KafkaTopicSpecArgs{
					Config:     convertmaps.ConvertGoMapToPulumiMap(config),
					Partitions: pulumi.Int(kafkaTopic.Partitions),
					Replicas:   pulumi.Int(kafkaTopic.Replicas),
					TopicName:  pulumi.String(kafkaTopic.Name),
				},
			}, pulumi.Parent(createdKafkaCluster))
		if err != nil {
			return errors.Wrapf(err, "failed to create kafka-topic %s", kafkaTopic.Id)
		}
	}
	return nil
}
