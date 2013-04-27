#!/usr/bin/python2.7
# Performs backup snapshots for ec2 hosts
# Writes status of backup to nagios
# if nagios isnt required, just change NAGIOS_CMD_FILE to whatever log
# this script depends on your amz keys being in the users home dir as described in boto docs
# add host to SNAPSHOT_HOSTS to start
import boto.ec2
import datetime
#takes Name from tags
SNAPSHOT_HOSTS = ['myhost1', 'myhost2']
EC2_REGION = 'us-east-1'
NAGIOS_CMD_FILE = '/var/spool/nagios/cmd/nagios.cmd'
#date,hostname, service, exitcode, output
NAGIOS_CHECK_TPL = """[%s] PROCESS_SERVICE_CHECK_RESULT;%s;%s;%s;%s"""
backupcount = {}

def connect_to_ec2():
    return boto.ec2.connect_to_region(EC2_REGION)

def get_all_reservations(conn):
    return conn.get_all_instances()

def get_all_instances(reservations):
    return [i for r in reservations for i in r.instances]

def connect_and_list():
    conn = connect_to_ec2()
    reservations = get_all_reservations(conn)
    instances = get_all_instances(reservations)
    return instances, conn


if __name__ == '__main__':
    backuphosts = {}
    now = datetime.datetime.now()
    epoch = now.strftime('%s')
    date = now.strftime('%Y-%d-%m')

    instances, conn = connect_and_list()

    for instance in instances:
        if instance.tags['Name'] in SNAPSHOT_HOSTS:
                    backuphosts[str(instance.id)] = instance.tags['Name']

    volumes = conn.get_all_volumes()

    for volume in volumes:
        vol = volume.attach_data
        
        if vol.instance_id in backuphosts:
            iid = vol.instance_id
            if iid in backupcount:
                backupcount[iid] += 1
            else:
                backupcount[iid] = 1

            snapname =  backuphosts[iid] +"-snapshot-"+str(backupcount[iid])+"-" +date
            nagios_submit = NAGIOS_CHECK_TPL % (epoch, "localhost", backuphosts[iid]+"-backup", 0, snapname )

            try:
                snapshot = conn.create_snapshot(vol.id, snapname)
                with open(NAGIOS_CMD_FILE, 'a') as file:
                    file.write(nagios_submit)
                
            except:
                nagios_submit = NAGIOS_CHECK_TPL % (epoch, "localhost", backuphosts[iid]+"-backup", 1, snapname )
                with open(NAGIOS_CMD_FILE, 'a') as file:
                    file.write(nagios_submit)
                     

