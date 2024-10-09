## tfunlock

Утилита для удаления записей из Serverless YDB для зависших или испорченых локов Terraform  
  
### Переменные окружения:  
AWS_ACCESS_KEY_ID="**AWSKEY**"  
AWS_SECRET_ACCESS_KEY="**AWSSECRET**"  
CUSTOM_ENDPOINT="https://docapi.serverless.yandexcloud.net/ru-central1/<путь к энпоинту Document API>"  
  
Параметры командной строки:  
-table название таблицы  

#### Пример использования:  
показать содержимое таблицы  
```shell
tfunlock -table prod-network-vpc-lock show
```
удалить запись о блокировке  
```shell
tfunlock -table prod-network-vpc-lock delete terraform-state-infra/dwh-test/net/terraform.tfstate-md5
```
поменять значение записи о блокиовке  
```shell
tfunlock -table prod-network-vpc-lock update terraform-state-infra/dwh-test/net/terraform.tfstate-md5 94bbd3501311e7569c1a350fd92b9162
```