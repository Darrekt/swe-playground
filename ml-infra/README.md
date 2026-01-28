1. Strategy & Architecture
  - Goal: Pivot from traditional DevOps (Jenkins, K8s, Bash) to ML Infrastructure (MLOps).
  - Platform: Chose Ubuntu 22.04 LTS as the host OS for its stability and broad ML library support.
  - Orchestrator: Selected kind (Kubernetes in Docker) over "from-scratch" installs to prioritize velocity and reproducibility (Cluster-as-Code).

2. Hardware & Drivers (The Host Layer)
  - GPU: AMD Radeon RX 7900 XTX.
  - Compute Stack: Installed ROCm 6.0.2 (Radeon Open Compute) using the amdgpu-install tool.
  - Decision - Hybrid Driver: Used the --no-dkms flag to layer ML libraries over the stable HWE kernel without replacing it.
  - Decision - Usecases: Installed graphics and rocm, but skipped mlsdk to minimize host-level bloat.
  - Permissions: Configured the render group to allow the user/docker to access the card.
  - Verification: Successful execution of rocm-smi on the host.

3. Orchestration Configuration (The Bridge)
  - Declarative Setup: Moving from imperative commands to a ml-cluster-config.yaml file.
  - Topology: Defined a 2-node cluster (Control-Plane + Worker) to simulate production environment isolation.
  - Hardware Passthrough: Mapped host devices (/dev/kfd, /dev/dri) into the Worker container.
  - Storage: Established a host-to-node mount (~/ml-data → /data) for persistent training data.

4. MLOps Advertisement (The Scheduler Layer)
  - Tool: Identified the AMD GPU Device Plugin for Kubernetes.
  - Mechanism: It will run as a DaemonSet to detect the hardware and update the Node's capacity with amd.com/gpu: 1.
  - Standardization: This allows the K8s scheduler to place pods based on GPU availability, just like it does for CPU/RAM.

