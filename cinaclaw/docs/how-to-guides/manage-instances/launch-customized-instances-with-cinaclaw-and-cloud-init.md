(how-to-guides-manage-instances-launch-customized-instances-with-cinaclaw-and-cloud-init)=
# Launch customized instances with cinaclaw and cloud-init

You can set up instances with a customized environment or configuration using the launch command along with a custom cloud-init YAML file and an optional post-launch health check to ensure everything is working correctly.

Below are some common examples of using `cloud-init` with cinaclaw to create customized instances. The `cloud-init` file is provided by the cinaclaw team, but users are free to create and use their own personal `cloud-init` configurations.

## 📦 anbox-cloud-appliance

Launch with:
```{code-block} text
cinaclaw launch \
  --name anbox-cloud-appliance \
  --cpus 4 \
  --memory 4G \
  --disk 50G \
  --timeout 900 \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-anbox.yaml
```

## ⚙️ charm-dev

Launch with:

```{code-block} text
cinaclaw launch 24.04 \
  --name charm-dev \
  --cpus 2 \
  --memory 4G \
  --disk 50G \
  --timeout 1800 \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-charm-dev.yaml
```

Health check:

```{code-block} text
cinaclaw exec charm-dev -- bash -c "
 set -e
 charmcraft version
 mkdir -p hello-world
 cd hello-world
 charmcraft init
 charmcraft pack
"
```

## 🐳 docker

Launch with:

```{code-block} text
cinaclaw launch 24.04 \
  --name docker \
  --cpus 2 \
  --memory 4G \
  --disk 40G \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-docker.yaml
```

Health check:

```{code-block} text
cinaclaw exec docker -- bash -c "docker run hello-world"
```

You can also optionally add aliases:

```{code-block} text
cinaclaw prefer docker
cinaclaw alias docker:docker docker
cinaclaw alias docker:docker-compose docker-compose
cinaclaw prefer default
cinaclaw aliases
```

> See also: [`How to use command aliases`](how-to-guides-manage-instances-use-instance-command-aliases)

## 🎞️ jellyfin

Launch with:

```{code-block} text
cinaclaw launch 22.04 \
  --name jellyfin \
  --cpus 2 \
  --memory 4G \
  --disk 40G \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-jellyfin.yaml
```

## ☸️ minikube

Launch with:

```{code-block} text
cinaclaw launch \
  --name minikube \
  --cpus 2 \
  --memory 4G \
  --disk 40G \
  --timeout 1800 \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-minikube.yaml
```

Health check:

```{code-block} text
cinaclaw exec minikube -- bash -c "set -e
  minikube status
  kubectl cluster-info"
```

## 🤖 ros2-humble

Launch with:

```{code-block} text
cinaclaw launch 22.04 \
  --name ros2-humble \
  --cpus 2 \
  --memory 4G \
  --disk 40G \
  --timeout 1800 \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-ros2-humble.yaml
```

Heath check:

```{code-block} text
cinaclaw exec ros2-humble -- bash -c "
  set -e

  colcon --help
  rosdep --version
  ls /etc/ros/rosdep/sources.list.d/20-default.list
  ls /home/ubuntu/.ros/rosdep/sources.cache

  ls /opt/ros/humble
"
```

## 🤖 ros2-jazzy

Launch with:

```{code-block} text
cinaclaw launch 24.04 \
  --name ros2-jazzy \
  --cpus 2 \
  --memory 4G \
  --disk 40G \
  --timeout 1800 \
  --cloud-init https://raw.githubusercontent.com/canonical/cinaclaw/refs/heads/main/data/cloud-init-yaml/cloud-init-ros2-jazzy.yaml
```

Health check:

```{code-block} text
cinaclaw exec ros2-jazzy -- bash -c "
  set -e

  colcon --help
  rosdep --version
  ls /etc/ros/rosdep/sources.list.d/20-default.list
  ls /home/ubuntu/.ros/rosdep/sources.cache

  ls /opt/ros/jazzy
"
```
