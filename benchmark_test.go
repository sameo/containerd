package containerd

import (
	"fmt"
	"testing"
)

func BenchmarkContainerCreate(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	client, err := New(address)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := testContext()
	defer cancel()

	image, err := client.GetImage(ctx, testImage)
	if err != nil {
		b.Error(err)
		return
	}
	spec, err := GenerateSpec(WithImageConfig(ctx, image), WithProcessArgs("true"))
	if err != nil {
		b.Error(err)
		return
	}
	var containers []Container
	defer func() {
		for _, c := range containers {
			if err := c.Delete(ctx); err != nil {
				b.Error(err)
			}
		}
	}()

	// reset the timer before creating containers
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("%s-%d", b.Name(), i)
		container, err := client.NewContainer(ctx, id, WithSpec(spec), WithNewRootFS(id, image))
		if err != nil {
			b.Error(err)
			return
		}
		containers = append(containers, container)
	}
	b.StopTimer()
}

func BenchmarkContainerStart(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	client, err := New(address)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := testContext()
	defer cancel()

	image, err := client.GetImage(ctx, testImage)
	if err != nil {
		b.Error(err)
		return
	}
	spec, err := GenerateSpec(WithImageConfig(ctx, image), WithProcessArgs("true"))
	if err != nil {
		b.Error(err)
		return
	}
	var containers []Container
	defer func() {
		for _, c := range containers {
			if err := c.Delete(ctx); err != nil {
				b.Error(err)
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("%s-%d", b.Name(), i)
		container, err := client.NewContainer(ctx, id, WithSpec(spec), WithNewRootFS(id, image))
		if err != nil {
			b.Error(err)
			return
		}
		containers = append(containers, container)

	}
	// reset the timer before starting tasks
	b.ResetTimer()
	for _, c := range containers {
		task, err := c.NewTask(ctx, empty())
		if err != nil {
			b.Error(err)
			return
		}
		defer task.Delete(ctx)
		if err := task.Start(ctx); err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}
