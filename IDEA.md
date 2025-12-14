# Idea

This project aims to create a simple tool to manage virtual machines (VMs) through a web interface. The idea is to make it easy for users to create, configure, start, stop, and delete VMs without dealing directly with Terraform or the underlying infrastructure.

# Concept

**Frontend:** A React-based interface where users can manage and monitor VMs.  
**Backend:** A GoLang API that handles requests from the frontend and uses Terraform to create and manage resources.  
**Backend/Provider:** At first, only Proxmox will be supported. Other providers could be added later.  
**Provisioning** (planned): Ansible could be used to automatically set up software and configurations after a VM is created.  

# Goal

The goal is to have a simple, centralized self-service portal for VMs, making it easier to manage infrastructure and reducing manual work for users.