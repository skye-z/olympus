# 待办事项

## Docker Registry

详见[官方文档](https://distribution.github.io/distribution/spec/api/)

* 上传镜像: PUT `/v2/<repository>/manifests/<tag>`
* 下载镜像: GET `/v2/<repository>/manifests/<tag>`
* 删除镜像: DELETE `/v2/<repository>/manifests/<tag>`
* 列出镜像: GET `/v2/<repository>/tags/list`
* 用户认证: POST `/v2/users/login`
* Token认证: POST `/v2/token/authenticate`

## Maven Registry

参考[GitLab文档](https://docs.gitlab.com/ee/api/packages/maven.html)

* [ ] 上传元数据: PUT `/<groupPath>/<artifactId>/<version>/maven-metadata.xml`
* [x] 下载元数据: GET `/<groupPath>/<artifactId>/<version>/maven-metadata.xml`
* [ ] 删除元数据: DELETE `/<groupPath>/<artifactId>/<version>/maven-metadata.xml`
* [ ] 上传制品: GET `/<groupPath>/<artifactId>/<version>/<artifactId>-<version>.<extension>`
* [x] 下载制品: GET `/<groupPath>/<artifactId>/<version>/<artifactId>-<version>.<extension>`
* [ ] 删除制品: DELETE `/<groupPath>/<artifactId>/<version>/<artifactId>-<version>.<extension>`
* [ ] 用户认证: POST `/users/login`
* [ ] Token认证: POST `/token/authenticate`

## NPM Registry

详见[官方文档](https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md)

* [ ] 包信息上传: PUT `/<package>`
* [x] 包信息下载: GET `/<package>`
* [ ] 特定版本包信息上传: PUT `/<package>/<version>`
* [x] 特定版本包信息下载: GET `/<package>/<version>`
* [ ] 包文件上传: PUT `/<package>/-/<filename>`
* [x] 包文件下载: GET `/<package>/-/<filename>`
* [ ] 用户认证: POST `/-/v1/login`
* [ ] Token认证: PUT `/-/v1/login`
