### hana-backup

This utility will perform a backup of an SAP HANA database and then copy the files off to Amazon S3.

## Requirements

* `~/.aws/credentials`  Your AWS credentials
* `~/.hdbsql`  This is an INI file and is documented below


## INI Definition

```
[hana]
host=<HANA HOST>
port=<HANA PORT>
user=<HANA USER>
password=<HANA PASSWORD>

[s3]
bucket=<DESTINATION BUCKET NAME>
region=<AWS REGION>

[backup]
location=<LOCAL BACKUP LOCATION>
prefix=<BACKUP PREFIX>
```

## Output
This utility will connect to the defined HANA instance and perform a backup with the following name:

`<PREFIX>-YYYY-M-D_H.M`

This is then uploaded to AWS S3


## Usage

```
[root@somehanaserver ~]# ./hana-backup-linux
BACKUP DATA USING FILE ('/backups/','SuperGreatApp-2017-7-19_15.44')
Uploading /backups/SuperGreatApp-2017-7-19_15.44_databackup_0_1 to us-east-1:SuperGreatApp-backups
Upload URL: https://SuperGreatApp-backups.s3.amazonaws.com/backups/SuperGreatApp-2017-7-19_15.44_databackup_0_1
Uploading /backups/SuperGreatApp-2017-7-19_15.44_databackup_2_1 to us-east-1:SuperGreatApp-backups
Upload URL: https://SuperGreatApp-backups.s3.amazonaws.com/backups%2FSuperGreatApp-2017-7-19_15.44_databackup_2_1
Uploading /backups/SuperGreatApp-2017-7-19_15.44_databackup_3_1 to us-east-1:SuperGreatApp-backups
Upload URL: https://SuperGreatApp-backups.s3.amazonaws.com/backups%2FSuperGreatApp-2017-7-19_15.44_databackup_3_1
```