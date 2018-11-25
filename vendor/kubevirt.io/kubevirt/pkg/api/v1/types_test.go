/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2017 Red Hat, Inc.
 *
 */

package v1

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/log"
)

var _ = Describe("PodSelectors", func() {
	Context("Pod affinity rules", func() {

		It("should work", func() {
			vmi := NewMinimalVMI("testvmi")
			vmi.Status.NodeName = "test-node"
			pod := &v1.Pod{}
			affinity := UpdateAntiAffinityFromVMINode(pod, vmi)
			newSelector := affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0]
			Expect(newSelector).ToNot(BeNil())
			Expect(len(newSelector.MatchExpressions)).To(Equal(1))
			Expect(len(newSelector.MatchExpressions[0].Values)).To(Equal(1))
			Expect(newSelector.MatchExpressions[0].Values[0]).To(Equal("test-node"))
		})

		It("should merge", func() {
			vmi := NewMinimalVMI("testvmi")
			vmi.Status.NodeName = "test-node"

			existingTerm := v1.NodeSelectorTerm{}
			secondExistingTerm := v1.NodeSelectorTerm{
				MatchExpressions: []v1.NodeSelectorRequirement{
					{},
				},
			}

			pod := &v1.Pod{
				Spec: v1.PodSpec{
					Affinity: &v1.Affinity{
						PodAffinity: &v1.PodAffinity{},
						NodeAffinity: &v1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
								NodeSelectorTerms: []v1.NodeSelectorTerm{
									existingTerm,
									secondExistingTerm,
								},
							},
						},
					},
				},
			}

			affinity := UpdateAntiAffinityFromVMINode(pod, vmi)

			Expect(affinity.NodeAffinity).ToNot(BeNil())
			Expect(affinity.PodAffinity).ToNot(BeNil())
			Expect(affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution).ToNot(BeNil())
			Expect(affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms).ToNot(BeNil())
			Expect(affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms).To(HaveLen(2))

			selector := affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0]
			Expect(selector).ToNot(BeNil())
			Expect(len(selector.MatchExpressions)).To(Equal(1))
			Expect(len(selector.MatchExpressions[0].Values)).To(Equal(1))
			Expect(selector.MatchExpressions[0].Values[0]).To(Equal("test-node"))

			selector = affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[1]
			Expect(selector).ToNot(BeNil())
			Expect(len(selector.MatchExpressions)).To(Equal(2))
			Expect(len(selector.MatchExpressions[1].Values)).To(Equal(1))
			Expect(selector.MatchExpressions[1].Values[0]).To(Equal("test-node"))
		})
	})
})

func TestSelectors(t *testing.T) {
	log.Log.SetIOWriter(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "PodSelectors")
}
