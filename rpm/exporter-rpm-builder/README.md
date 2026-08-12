# Cluster Exporter RPM Build Tool

Creates a virtual linux machine, then builds the RPM packages for the Prometheus Cluster Exporter using it. Works with version 1.1.8 onwards.

## Requirements

Vagrant, Ansible and Virtualbox or Libvirt. 

## Instructions

Edit `Vagrantfile` to choose your linux version (Fedora 43 or Rocky 8) and your vm provider (Libvirt or Virtualbox). Edit `group_vars/all.yml` to set version number, command line options and the URL for the Cluster Exporter git repository. Run `vagrant up`.

Generated RPMs are found in the `rpms/` subfolder.

