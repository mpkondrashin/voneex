# VOneEx - Avoid adding Special-Purpose addresses to Vision One Suspicious Objects list

**Populate Vision One IoC exception list with IANA IPv4 Special-Purpose Address**

<p align="center">
  <img src="voneex.jpeg" alt="VOneEx" width="50%">
</p>


VOneEx populates the Vision One IoC exception list with [IANA IPv4 Special-Purpose Address Registry](https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml) entries to prevent false positives and unnecessary alerts.

## Usage

### Get Vision One Token
1. Login to [Vision One Portal](https://portal.xdr.trendmicro.com/index.html)
2. Create limited User Role
    1. Go to Administration -> User Roles
    2. Press "Add Role" button
    3. Role name: "VOneEx"
    4. On permissions tab->Threat Intelligence->Suspicious Object Management - check all checkboxes 
    5. Press "Save" button

3. Create API Key
    1. Go to Administration->API Keys
    2. Press "Add API Key" button
    3. Name: "VOneEx", Role: "VOneEx" created previously
    4. Press "Add" button
    5. Save API Key for later use

### Get The Latest Version

Download [latest release](https://github.com/mpkondrashin/voneex/releases/latest) for your platform and unpack archive.

### Run VOneEx

Run following command:
```commandline
voneex --apikey <api key>
```

## Bugs:

### If some of IANA Address Registry IPs are not needed

Go to Threat Intelligence->Suspicious Object Management->Exception List->Object type->IP address and delete excess entries.

### If macOS does not allow to run voneex, run following command:
```commandline
xattr -c voneex
```