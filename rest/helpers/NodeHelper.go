package helpers

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"strconv"
)

//第一个是cpu使用 第二个是内存使用
func GetNodeUsage(c *versioned.Clientset, node *v1.Node) (cpu, mem float64) {
	nodeMetric, _ := c.MetricsV1beta1().
		NodeMetricses().Get(context.Background(), node.Name, metav1.GetOptions{})
	cpu = float64(nodeMetric.Usage.Cpu().MilliValue()) / float64(node.Status.Capacity.Cpu().MilliValue())
	mem = float64(nodeMetric.Usage.Memory().MilliValue()) / float64(node.Status.Capacity.Memory().MilliValue())
	return Decimal(cpu), Decimal(mem)
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}
