# Local environment setpu

### IDE

- VSCode

- Install recommended extensions for VS Code (As per [/.vscode/extensions.json](/.vscode/extensions.json))

### System setup

- git

- Use WSL2 for windows: [Ubuntu subsystem: 24.04.2](https://ubuntu.com/desktop/wsl)

    ```ps1
    wsl --set-default-version 2
    wsl --install -d Ubuntu-24.04
    # Offline installation: Get file from https://ubuntu.com/desktop/wsl
    wsl --install --from-file D:\apps\ubuntu-24.04.2-wsl-amd64.wsl
    wsl --list --all
    wsl -d Ubuntu-24.04 --shutdown
    # Optional: if you do not wish to use default disk location, e.g. let's say on D:
    wsl --export Ubuntu-24.04 D:\wsl-ubuntu-24.04.tar
    wsl --unregister Ubuntu-24.04
    wsl --import Ubuntu-24.04 "D:\wsl2\ubuntu-24.04" "D:\wsl-ubuntu-24.04.tar"
    wsl --set-default Ubuntu-24.04
    ```

- docker >= 28.1.1 (or docker desktop)

    ```sh
    # For windows install cli and install docker engine on WSL
    choco install docker-cli
    
    # Install docker engine
    # Reference: https://docs.docker.com/engine/install/ubuntu/
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc
    # Add the repository to Apt sources:
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update

    # Use docker without sudo
    sudo groupadd docker
    sudo gpasswd -a $USER docker
    sudo usermod -a -G docker $USER
    chmod 664 /var/run/docker.sock
    grep docker /etc/group
    # Exit from wsl and login again
    ```

    Set docker to listen on both file and tcp socket

    Get dockerd service configuration path from following command

    ```sh
    sudo service docker status
    sudo nano /usr/lib/systemd/system/docker.service
    ```

    Add `-H tcp://127.0.0.1:2375` in `ExecStart`, e.g.

    ```
    ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2375 --containerd=/run/containerd/containerd.sock
    ```

    Reload daemon and restart docker service

    ```sh
    sudo systemctl daemon-reload
    sudo service docker restart
    ```

    Set docker context on windows host

    Note: If you don't wish to set docker context, you can use using

    1. Environment variable

    ```sh
    export DOCKER_HOST=tcp://127.0.0.1:2375
    ```

    2. Use `--host` parameter

    ```sh
    docker --host tcp://127.0.0.1:2375 info
    ```

    ```sh
    # Create context named "wsl-ubuntu" that uses WSL Ubuntu as host
    docker context create wsl-ubuntu --docker "host=tcp://127.0.0.1:2375"
    # Set wsl as default context
    docker context use wsl-ubuntu
    # Check if wsl is default context
    docker context list
    ```

    You'll see above changes in `~/.docker/config.json` & `~/.docker/daemon.json`

- kubernetes ~ v1.31

    - [k3d](https://github.com/k3d-io/k3d) >= v5.8.3 + k3s >= v1.31.5-k3s1 or [kind](https://github.com/kubernetes-sigs/kind)

        For k3d+k3s cluster: refer [./k3s-cluster-config.yaml](./k3s-cluster-config.yaml).

        Refer [k3d config doc](https://k3d.io/v5.8.3/usage/configfile)

        Note: K3s uses `rancher/k3s` image about `222MB` to create nodes

        ```sh
        k3d cluster create mycluster --config ./k3s-cluster-config.yaml

        # You can change load-balancer poart using
        k3d cluster edit mycluster --port-add 8085:80

        # To delete cluster
        k3d cluster delete mycluster
        ```

        Refer [kind cluster docs](https://kind.sigs.k8s.io/docs/user/quick-start/)

        Note: Kind uses `kindest/node` image about `1.07GB` to create nodes

        ```sh
        kind create cluster --help

        kind create cluster # Default cluster context name is `kind`.
        # kind create cluster --name kind-2

        kind get clusters

        # Delete kind cluster
        kind delete cluster # Default cluster context name is `kind`.
        # kind delete cluster --name kind-2

        # Docker images can be loaded into your cluster nodes with:
        kind load docker-image my-custom-image-0 my-custom-image-1
        ```

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) ~= v1.33.1

    For WSL, create softlink in WSL ubuntu to windows host so that both WSL and Windows host connect to same kubecontext, e.g.

    ```sh
    ln -s /mnt/c/Users/<Your windows USERNAME here>/ ~/.kube
    ```

- [helm](https://github.com/helm/helm) >= v3.18.0-rc.2

- [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) >= 4.5.2

- [asciinema](https://github.com/asciinema/asciinema/releases/tag/v3.0.0-rc.4) >= 3.0.0-rc.4

- [Golang](https://go.dev/dl/) >= 1.24.3

    ```sh
    # For linux or WSL
    wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
    tar -xvzf go1.24.3.linux-amd64.tar.gz
    mv /usr/lib/go-1.24.3
    sudo ln -s /usr/lib/go-1.24.3 /usr/lib/go
    sudo ln -s /usr/lib/go-1.24.3/bin/go /usr/bin/go
    sudo ln -s /usr/lib/go-1.24.3/bin/gofmt /usr/bin/gofmt
    ```

    Check go installation

    ```sh
    go version
    go env
    ```

- make >= 4.3

    ```sh
    # For linux or WSL
    sudo apt-get install make
    ```