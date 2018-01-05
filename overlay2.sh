#!/usr/bin/env bash
#overlay2 modules
overlayKo="/lib/modules/3.10.0-327.44.58.35.x86_64/kernel/fs/overlayfs/overlay.ko"
modulesLoad="/etc/modules-load.d"
dockerOverlay="${modulesLoad}/overlay.conf"

#docker daemon driver
dockerConfDir="/etc/docker"
dockerDaemon="${dockerConfDir}/daemon.json"

#docker proxy
dockerServiceDir="./etc/systemd/system/docker.service.d"
dockerProxy="${dockerServiceDir}/http-proxy.conf"

#docker mirror
dockerServiceLib="/lib/systemd/system/docker.service"
dockerServiceEtc="/etc/systemd/system/docker.service"
execStart="ExecStart=/usr/bin/docker daemon"
registry_mirror="--registry-mirror=https://100.125.5.235:20202"


needUpdateOverlay=0
if [ ! -d "${modulesLoad}" ]; then
    mkdir -p  ${modulesLoad}
fi

if [ -f "${dockerOverlay}" ]; then
    needUpdateOverlay=`cat "${dockerOverlay}" | grep 'overlay' | wc -l`
fi

if [ ${needUpdateOverlay} -eq 0 ]; then
    echo "overlay" >> ${dockerOverlay}
fi
#insmod overlay driver
insmod ${overlayKo}

needUpdateDaemon=0
if [ ! -d "${dockerConfDir}" ]; then
    mkdir -p  ${dockerConfDir}
fi

if [ -f "${dockerDaemon}" ]; then
   needUpdateDaemon=`cat "${dockerDaemon}" | grep '"storage-driver": "overlay2"' | wc -l`
fi

if [ ${needUpdateDaemon} -eq 0 ]; then
    #set docker storage driver
    echo '{' >> ${dockerDaemon}
    echo '  "storage-driver": "overlay2",' >> ${dockerDaemon}
    echo '  "storage-opts": [' >> ${dockerDaemon}
    echo '      "overlay2.override_kernel_check=true"' >> ${dockerDaemon}
    echo '  ]' >> ${dockerDaemon}
    echo '}' >> ${dockerDaemon}
fi

needUpdateService=0
if [ ! -d "${dockerServiceDir}" ]; then
    mkdir -p  ${dockerServiceDir}
fi

if [ -f "${dockerProxy}" ]; then
   needUpdateService=`cat "${dockerProxy}" | grep '10.177.221.12:3128' | wc -l`
fi

if [ ${needUpdateService} -eq 0 ]; then
    #set docker proxy
    echo '[Service]' >> ${dockerProxy}
    echo 'Environment="HTTP_PROXY=10.177.221.12:3128" "HTTPS_PROXY=10.177.221.12:3128" "NO_PROXY=100.125.5.235"' >> ${dockerProxy}
fi

needUpdateMirror=0
if [ -f "${dockerServiceEtc}" ]; then
   needUpdateMirror=`cat "${dockerServiceEtc}" | grep "${registry_mirror}" | wc -l`
else
   needUpdateMirror=`cat "${dockerServiceLib}" | grep "${registry_mirror}" | wc -l`
   if [ ${needUpdateMirror} -eq 0 ]; then
        cp -n ${dockerServiceLib} ${dockerServiceEtc}
   fi
fi

if [ ${needUpdateMirror} -eq 0 ]; then
    sed -i "s|${execStart}|${execStart} ${registry_mirror}" ${dockerServiceEtc}
fi

#restat docker daemon
systemctl daemon-reload
systemctl restart docker
