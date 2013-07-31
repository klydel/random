#!/usr/bin/python2.7
from git import *
repodir = '/home/klydel/git'

class AwsInfo():
    def __init__(self):
        self.status = {}
    
    def get_ec2_info(self):
        import boto.ec2
        conn = boto.ec2.connect_to_region('us-east-1')
        reservations = conn.get_all_instances()
        for reserv in reservations:
        
            for serv in reserv.instances:
              
                self.status[str(reserv.id)] = {serv : {'instance_type' : serv.instance_type, 'groups' : [str(x.name) for x in serv.groups], 'key_name' : serv.key_name, 'image_id' : serv.image_id, 'placement' : serv.placement, 'kernel' : serv.kernel, 'ramdisk' : serv.ramdisk, 'architecture' : serv.architecture, 'block_device_mapping' : [{x : y.volume_id} for x,y in serv.block_device_mapping.iteritems()], 'root_device_name' : serv.root_device_name, 'root_device_type' : serv.root_device_type, 'ebs_optimized' : serv.ebs_optimized, 'tags' : [{k : v} for k,v in serv.tags.iteritems()]}}
              
        return self.status
        

class Route53Info():
    def __init__(self):
        pass

    def get_r53_info(self):
        import boto.route53
        conn = boto.route53.connection.Route53Connection()
        self.zones = dict(conn.get_all_hosted_zones())
        self.Ids = [x.Id.replace('/hostedzone/', '') for x in self.zones['ListHostedZonesResponse']['HostedZones']]
        self.rr = [conn.get_all_rrsets(str(x)) for x in self.Ids]
        return self.rr



def write_git_info(astatus, zstatus):

    repo = Repo(repodir)
    origin = repo.remotes.origin
    origin.pull()
    with open(repodir + 'ec2-inventory', 'w') as file:
        file.write(str(astatus))
    with open(repodir + 'r53-inventory', 'w') as file:
        file.write(str(zstatus))
    

    index = repo.index
    for d in repo.index.diff(None):
        if 'inventory' in str(d):
            print d.a_blob.name
            index.add(['info/'+d.a_blob.name])
            index.commit("Production Amazon Web Services Change: Updated Inventory Files") 
            index.write()
    origin.push()
    


if __name__ == '__main__':
    a = AwsInfo()
    astatus = a.get_ec2_info()
    z = Route53Info()
    zstatus = z.get_r53_info()
    write_git_info(astatus, zstatus)
