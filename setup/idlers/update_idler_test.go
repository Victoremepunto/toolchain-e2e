package idlers

import (
	"context"
	"testing"
	"time"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	commontest "github.com/codeready-toolchain/toolchain-common/pkg/test"
	cfg "github.com/codeready-toolchain/toolchain-e2e/setup/configuration"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestUpdateTimeout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cfg.DefaultTimeout = time.Second * 1
		idler := &toolchainv1alpha1.Idler{
			ObjectMeta: metav1.ObjectMeta{
				Name: "user0001-dev",
			},
			Spec: toolchainv1alpha1.IdlerSpec{
				TimeoutSeconds: 43200,
			},
			Status: toolchainv1alpha1.IdlerStatus{
				Conditions: []toolchainv1alpha1.Condition{
					{
						Type:   toolchainv1alpha1.ConditionReady,
						Status: corev1.ConditionTrue,
						Reason: "Running",
					},
				},
			},
		}
		cl := commontest.NewFakeClient(t, idler)

		// when
		err := UpdateTimeout(cl, "user0001", 15*time.Second)

		// then
		require.NoError(t, err)
		updatedIdler := &toolchainv1alpha1.Idler{}
		require.NoError(t, cl.Get(context.TODO(), types.NamespacedName{Name: "user0001-dev"}, updatedIdler))
		assert.Equal(t, int32(15), updatedIdler.Spec.TimeoutSeconds)
	})

	t.Run("compliant username differs from usersignup name", func(t *testing.T) {
		// given
		// Reproduces issue #1308: when the sandbox operator transforms
		// the username (e.g. "default-0005" -> "crt-default-0005"),
		// the Idler is named "crt-default-0005-dev", not "default-0005-dev".
		cfg.DefaultTimeout = time.Millisecond * 100
		idler := &toolchainv1alpha1.Idler{
			ObjectMeta: metav1.ObjectMeta{
				Name: "crt-default-0005-dev",
			},
			Spec: toolchainv1alpha1.IdlerSpec{
				TimeoutSeconds: 43200,
			},
			Status: toolchainv1alpha1.IdlerStatus{
				Conditions: []toolchainv1alpha1.Condition{
					{
						Type:   toolchainv1alpha1.ConditionReady,
						Status: corev1.ConditionTrue,
						Reason: "Running",
					},
				},
			},
		}
		cl := commontest.NewFakeClient(t, idler)

		// when - using the original username (not the compliant one)
		err := UpdateTimeout(cl, "default-0005", 15*time.Second)

		// then - should find the idler even when the name was transformed
		require.NoError(t, err)
		updatedIdler := &toolchainv1alpha1.Idler{}
		require.NoError(t, cl.Get(context.TODO(), types.NamespacedName{Name: "crt-default-0005-dev"}, updatedIdler))
		assert.Equal(t, int32(15), updatedIdler.Spec.TimeoutSeconds)
	})
}
