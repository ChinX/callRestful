#!/usr/bin/env bash
#镜像仓库地址
registryAdd="100.125.5.235:20202"

#运行场景: 类生产为"prepare"; 生产环境为"formally"
runtimeScene="prepare"

#euler内核版本：一般不需修改
eulerOSKernel="3.10.0-327.44.58.35.x86_64"

echo '
######################关于配置项的说明######################

镜像仓库：     registryAdd="'${registryAdd}'"
运行场景：     runtimeScene="'${runtimeScene}'"    ###类生产为"prepare"; 生产环境为"formally"
euler内核版本：eulerOSKernel="'${eulerOSKernel}'"  ###一般不需要配置

可通过命令行参数传入: bash overlay2.sh {registryAdd} {runtimeScene} {eulerOSKernel}
参数可由后至前依次省略，省略项则使用以上对应的默认值

##########################################################'

if [ $# -ge 1 ]; then
    if [[ $1 == "-h" || $1 == "--help" ]]; then
       exit 0
    fi
    registryAdd=$1
    if [ $# -ge 2 ]; then
        if [ $2 == "formally" ]; then
            runtimeScene="formally"
        fi
        if [ $# -ge 3 ]; then
            eulerOSKernel=$3
        fi
    fi
fi

#overlay2 modules
overlayKo="/lib/modules/${eulerOSKernel}/kernel/fs/overlayfs/overlay.ko"
modulesLoad="/etc/modules-load.d"
dockerOverlay="${modulesLoad}/overlay.conf"

#docker daemon driver
dockerConfDir="/etc/docker"
dockerDaemon="${dockerConfDir}/daemon.json"

#docker proxy
dockerServiceDir="/etc/systemd/system/docker.service.d"
dockerProxy="${dockerServiceDir}/http-proxy.conf"

#docker mirror
dockerServiceLib="/lib/systemd/system/docker.service"
dockerServiceEtc="/etc/systemd/system/docker.service"


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

#生产环境pod域直通外网，不需要设置代理
if [ ${runtimeScene} != "formally" ]; then
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
fi

needUpdateMirror=0
if [ -f "${dockerServiceEtc}" ]; then
   needUpdateMirror=`cat "${dockerServiceEtc}" | grep "registry-mirror" | grep "${registryAdd}" | wc -l`
else
   needUpdateMirror=`cat "${dockerServiceLib}" | grep "registry-mirror" | grep "${registryAdd}" | wc -l`
   if [ ${needUpdateMirror} -eq 0 ]; then
        cp -n ${dockerServiceLib} ${dockerServiceEtc}
   fi
fi

if [ ${needUpdateMirror} -eq 0 ]; then
    sed -i 's/docker daemon/docker daemon --registry-mirror=https:\/\/'${registryAdd}'/g' ${dockerServiceEtc}
fi

#restat docker daemon
systemctl daemon-reload
systemctl restart docker
