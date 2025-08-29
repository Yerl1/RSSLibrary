package workerpool

type WorkerLauncher interface {
	LaunchWorker(in chan string, stopCh chan struct{})
}
