# 构建docker iamge
tag=0.12

games=(connector game1)

for game in ${games[@]}
do 
    echo -e "\n=====================================================================\n"

    cd $game; echo `pwd`

    echo 正在构建${game}镜像;

    rm -rf vendor # 删除之前的依赖包

    go get -u # 更新依赖

    go mod vendor # 将依赖放在项目内，统一打包到docker镜像,后面就可以不用安装直接使用了

    docker build . -t $game # build镜像

    docker tag $game 127.0.0.1:5000/${game}:${tag} # 打标签

    docker push 127.0.0.1:5000/${game}:${tag}  # 存放到本地私有仓库

    cd ../
done