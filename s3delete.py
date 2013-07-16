#!/usr/bin/python2.7
#simple bulk delete script, use at your own risk
#use keepfiles list for files you dont want to delete
#use keepdir for a more generic s3 filename keeplist
#change TEST to false to actually delete
#change BOTO_CFG to whereever your amz keys are
from boto.s3.connection import S3Connection
TEST = True
keepdir = ['storagebackup']
keepfiles = ['backup/git2013-07-03.gz', 'backup/git2013-07-04.gz']
BOTO_CFG = '.boto.cfg'
S3_BUCKET = 'mybucket'

def parse_config():
   import ConfigParser
   config = ConfigParser.ConfigParser()
   config.read([BOTO_CFG])
   aws_access_key_id = config.get('Credentials', 'aws_access_key_id')
   aws_secret_access_key = config.get('Credentials', 'aws_secret_access_key')
   return aws_access_key_id,aws_secret_access_key



if __name__ == '__main__':

    aws_access, aws_secret = parse_config()
    conn = S3Connection(aws_access, aws_secret)
    b = conn.get_bucket(S3_BUCKET)

    for key in b.list():
        if key.name not in keepfiles and key.name.split('/')[0] not in keepdir:
            print "DELETING", key.name
            if not TEST:
               key.delete()
        else:
            print "KEEPING", key.name
