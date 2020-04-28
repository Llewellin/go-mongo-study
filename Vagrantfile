# -*- mode: ruby -*-
# vi:set ft=ruby sw=2 ts=2 sts=2:

# Define the number of master and worker nodes
# If this number is changed, remember to update setup-hosts.sh script with the new hosts IP details in /etc/hosts of each VM.
NUM_MASTER_NODE = 1
NUM_WORKER_NODE = 2

IP_NW = "44.44.44."
MASTER_IP_START = 10
NODE_IP_START = 20
LB_IP_START = 30

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  # config.vm.box = "base"
  config.vm.box = "ubuntu/bionic64"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  config.vm.box_check_update = false

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  

  # Provision Worker Nodes
  (1..NUM_WORKER_NODE).each do |i|
    config.vm.define "mongo-secondary-#{i}" do |node|
        node.vm.provider "virtualbox" do |vb|
            vb.name = "mongo-secondary-#{i}"
            vb.memory = 4096
            vb.cpus = 2
        end
        node.vm.hostname = "mongo-secondary-#{i}"
        node.vm.network :private_network, ip: IP_NW + "#{NODE_IP_START + i}"
		# node.vm.network "forwarded_port", guest: 22, host: "#{2720 + i}"

        node.vm.provision "setup-hosts", :type => "shell", :path => "ubuntu/setup-hosts.sh" do |s|
          s.args = ["enp0s8"]
        end

        node.vm.provision "install-mongo", type: "shell", :path => "ubuntu/install-mongo.sh"
        node.vm.provision "start-mongo#{i+1}", type: "shell", :path => "ubuntu/start-mongo#{i+1}.sh"

        # node.vm.provision "setup-dns", type: "shell", :path => "ubuntu/update-dns.sh"
        # node.vm.provision "install-docker", type: "shell", :path => "ubuntu/install-docker3.sh"
        # node.vm.provision "allow-bridge-nf-traffic", :type => "shell", :path => "ubuntu/allow-bridge-nf-traffic.sh"
        # node.vm.provision "install-kubeadm", :type => "shell", :path => "ubuntu/install-kubeadm.sh"
    end
  end

  # Provision Master Nodes
  (1..NUM_MASTER_NODE).each do |i|
    config.vm.define "mongo-primary-#{i}" do |node|
        # Name shown in the GUI
        node.vm.provider "virtualbox" do |vb|
            vb.name = "mongo-primary-#{i}"
            vb.memory = 4096
            vb.cpus = 2
        end
        node.vm.hostname = "mongo-primary-#{i}"
        node.vm.network :private_network, ip: IP_NW + "#{MASTER_IP_START + i}"
        # node.vm.network "forwarded_port", guest: 22, host: "#{2710 + i}"
        # node.vm.network "forwarded_port", guest: 30008, host: "30008"

        node.vm.provision "setup-hosts", :type => "shell", :path => "ubuntu/setup-hosts.sh" do |s|
          s.args = ["enp0s8"]
        end

        node.vm.provision "install-mongo", type: "shell", :path => "ubuntu/install-mongo.sh"
        node.vm.provision "start-mongo", type: "shell", :path => "ubuntu/start-mongo.sh"
        node.vm.provision "setup-rs", type: "shell", :path => "ubuntu/setup-rs.sh"
        # # customized
        # node.vm.provision "install-docker", type: "shell", :path => "ubuntu/install-docker3.sh"
        # node.vm.provision "install-kubeadm", :type => "shell", :path => "ubuntu/install-kubeadm.sh"
        # node.vm.provision "init-kubeadm", :type => "shell", :path => "ubuntu/init-kubeadm.sh"
        # node.vm.provision "copy-initial-files", :type => "shell", :path => "ubuntu/copy-initial-files.sh"
        # node.vm.provision "install-ingress", :type => "shell", :path => "ubuntu/install-ingress.sh"
        # node.vm.provision "install-tidal", :type => "shell", :path => "ubuntu/install-tidal.sh"
    end
  end
end
