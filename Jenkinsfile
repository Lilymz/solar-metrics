pipeline {
    agent any
    environment{
       URL = "https://github.com/Lilymz/solar-metrics.git"
       GIT_CRED = "solar"
       APP_NAME = "solar-metrics"
       VERSION = "V1.0.0"
    }
    stages{
        stage("Dockerfile Before Handler"){
           steps{
                script{
                    // 清理存在的镜像，
                    sh '''
                        docker rmi "${APP_NAME}"
                        docker rm -f "${appName}"
                        docker image prune -f
                        docker network prune -f
                    '''
                }
           }
        }
        stage("Checkout"){
            steps{
                script{
                    checkout([
                     $class: 'GitSCM',
                     branches: [[name: '*/develop']],
                     doGenerateSubmoduleConfigurations: false, extensions: [],
                     userRemoteConfigs: [[
                        url: "${URL}",
                        credentialsId: "${GIT_CRED}"
                       ]]
                     ])
                }
            }
        }
        stage("Build"){
            steps{
                script{
                    sh '''
                        go mod tidy
                        sleep 3
                        go build -o "${APP_NAME}" internal/.
                    '''
                }
            }
        }
        stage("Build Image"){
            steps{
                script{
                    // 构建镜像
                    sh(script:"docker build --no-cache --compress --pull --tag ${APP_NAME}:${VERSION} .")
                }
            }
        }
        stage("Deploy"){
            steps{
                script{
                    sh '''
                        docker compose up -d
                    '''
                }
            }
        }
    }
    post{
        success{
            script{
                // 构建成功，将当前镜像推送到私服仓库
                // 发送构建成功通知
            }
        }
    }
}
