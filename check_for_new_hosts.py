#!/usr/bin/python2.7
#crappy attempt to add hosts to nagios and route53 via cron. creates individual host confs and r53 records, then sends email confirmation.
# layout in NAGIOS_CONF_PATH uses the value from NAGIOS_CONF_MAPPING.  NAGIOS_CONF_MAPPING value is used as a tag in ec2 tags feature.
#NAGIOD_CONF_MAPPING key is also used for servicegroups, hostgroups.
#NAGIOS_CONF_MAPPING maps one or many tags to single directory layout.
#only adds hosts defined with an ec2 tag of EC2_TAG, skips all other hosts.
#also skips hosts that are in NAGIOS_HOST_SKIP.
#adds dns record using ZONE_ID.
#EMAIL_ALERT set to True sends email via localhost alerting that nagios is ready to be restarted with new confs.
#use in a test env first, use at your own risk, etc, etc.
#klydel@gmail.com

import boto
import boto.ec2
import os

EMAIL_ALERT = True
BODY = ""
EMAIL_TO = ["oncall@company.com"]
EMAIL_FROM = 'nagios@company.com'
HOSTNAME_POSTFIX = '.company.com'
EC2_REGION = "us-east-1"
EC2_TAG = 'serverrole'
ZONE_ID = 'xxxxxxxxxx'
NAGIOS_CONF_PATH = '/etc/nagios/'
NAGIOS_HOST_SKIP = ['qa','dev']
NAGIOS_CONF_MAPPING = {
'production-www':'production-www',
'prod-application':'prod-app',
'prod-celery-worker':'prod-celery',
'prod-celery-worker,prod-celery-worker1,prod-celery-worker2':'prod-celery',
'prod-mongo-data':'prod-mongo'
}
NAGIOS_HOST_TPL = """
define host{
        use                     linux-server 
        host_name               %s
        hostgroups              %s
        alias                   %s
        address                 %s
        }
"""

def check_host_conf(hosttag, host):
    try:
        return os.path.isfile(NAGIOS_CONF_PATH + NAGIOS_CONF_MAPPING[hosttag] +"/"+host+".cfg")
    except:
        return False

def connect_to_ec2():
    return boto.ec2.connect_to_region(EC2_REGION)

def get_all_reservations(conn):
    return conn.get_all_instances()

def get_all_instances(reservations):
    return [i for r in reservations for i in r.instances]

def get_all_tags(reservations):
    return [[x.tags for x in reservations[i].instances] for i in xrange(len(reservations))]

def get_instances_roles(alltags):
    return [i for i in alltags if EC2_TAG in i[0]]

def get_instance_roles(instances):
    return [[i.tags, i.private_dns_name] for i in instances if EC2_TAG in i.tags if i.state == 'running']

def write_template(hostname, filecontents, role):
    try:
        hostcfgfile = NAGIOS_CONF_PATH + NAGIOS_CONF_MAPPING[role] +"/"+hostname+".cfg"
        with open(hostcfgfile, "w") as nagiosconf:
            nagiosconf.write(filecontents)
    except IOError:
        print "Unable to write configuration file"

def check_config():
    pass

def restart_nagios():
    pass

def send_alert(body):
        if EMAIL_ALERT == True:
          try:
              import smtplib
              from email.mime.text import MIMEText
              msg = MIMEText(body)
              msg['Subject'] = 'Added New Host to Nagios'
              msg['From'] = EMAIL_FROM
              msg['To'] = ", ".join(EMAIL_TO)
              s = smtplib.SMTP('localhost')
              s.sendmail(EMAIL_FROM, EMAIL_TO, msg.as_string())
              s.quit()

          except:
              print "alert failed"

def connect_and_list():
    conn = connect_to_ec2()
    reservations = get_all_reservations(conn)
    instances = get_all_instances(reservations)
    instanceroles = get_instance_roles(instances)
    return instanceroles

def create_dns_record(hostname, hostip):
    from boto.route53.record import ResourceRecordSets
    conn = boto.connect_route53()
    edits = ResourceRecordSets(conn, ZONE_ID)
    edit = edits.add_change("CREATE", hostname.strip() + HOSTNAME_POSTFIX, "CNAME")
    edit.add_value(hostip)
    status = edits.commit()

if __name__ == '__main__':
    instanceroles = connect_and_list()
    for i in instanceroles:
        if i[0].get('Name') in NAGIOS_HOST_SKIP:
            pass
        else:
            #fix
            res = check_host_conf(str(i[0].get(EC2_TAG)), str(i[0].get('Name')))
            if not res:
                try:
                    hostname = str(i[0].get('Name'))
                    servicegroups = i[0].get(EC2_TAG) + ","+NAGIOS_CONF_MAPPING[i[0].get(EC2_TAG)]+"-base"
                    print "Found new host: ", i[0].get('Name'), i[0].get(EC2_TAG)
                    
                    filecontents = NAGIOS_HOST_TPL % (hostname, servicegroups, hostname, hostname + HOSTNAME_POSTFIX)
                    write_template(hostname, filecontents, i[0].get(EC2_TAG))
                    try:
                        create_dns_record(i[0].get('Name'), i[1])
                    except:
                        print "failed to create DNS records"
                    if EMAIL_ALERT:
                        BODY += "Host: %s , Role: %s \r\n" % (hostname, i[0].get(EC2_TAG))
                except:
                    print "skipping non roled host"
                    pass

    if EMAIL_ALERT:
        try:
            if len(BODY) > 1:
                send_alert(BODY)
        except:
            pass
    
                    

                
