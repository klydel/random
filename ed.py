#!/usr/bin/python

import sys
import subprocess, signal

WARE_MAPPING = {
    'remvsearch' : {
        'order' : ['files'],
        'files' : ['/Library/Application Support/VSearch', '/Library/LaunchAgents/com.vsearch.agent.plist', '/Library/LaunchDaemons/com.vsearch.daemon.plist', '/Library/LaunchDaemons/com.vsearch.helper.plist', '/Library/LaunchDaemons/Jack.plist', '/Library/PrivilegedHelperTools/Jack', '/System/Library/Frameworks/VSearch.framework']
    },
    'remconduit' : {
        'order' : ['files'],
        'files' : ['/Library/InputManagers/CTLoader/', '/Library/LaunchAgents/com.conduit.loader.agent.plist', '/Library/LaunchDaemons/com.perion.searchprotectd.plist', '/Library/Application Support/SIMBL/Plugins/CT2285220.bundle', '/Library/Application Support/Conduit/', '/Applications/SearchProtect.app', '~/Library/Application Support/Conduit/', '~/Library/Internet Plug-Ins/ConduitNPAPIPlugin.plugin', '~/Library/Internet Plug-Ins/TroviNPAPIPlugin.plugin', '~/Conduit/', '~/Trovi/']
    },
    'installmac' : {
        'order' : ['kill', 'files'],
        'kill' : ['Genieo', 'InstallMac'],
        'files' : ['/private/etc/launchd.conf', '/Applications/Genieo', '/Applications/InstallMac', '/Applications/Uninstall Genieo<', '/Applications/Uninstall IM Completer.app', '~/Library/Application Support/com.genieoinnovation.Installer/', '~/Library/Application Support/Genieo/', '~/Library/LaunchAgents/com.genieo.completer.download.plist', '~/Library/LaunchAgents/com.genieo.completer.update.plist', '/Library/LaunchAgents/com.genieoinnovation.macextension.plist', '/Library/LaunchAgents/com.genieoinnovation.macextension.client.plist', '/Library/LaunchAgents/com.genieo.engine.plist', '/Library/LaunchAgents/com.genieo.completer.update.plist', '/Library/LaunchDaemons/com.genieoinnovation.macextension.client.plist', '/Library/PrivilegedHelperTools/com.genieoinnovation.macextension.client', '/usr/lib/libgenkit.dylib', '/usr/lib/libgenkitsa.dylib', '/usr/lib/libimckit.dylib', '/usr/lib/libimckitsa.dylib', '/Library/Frameworks/GenieoExtra.framework']
    },
}

Test = False

def get_process_list():
    p = subprocess.Popen(['ps', '-A'], stdout=subprocess.PIPE)
    out, err = p.communicate()
    return out

def kill_app(app):
    proc_list = get_process_list()
    for line in out.splitlines():
        if app in line:
            try:
                pid = int(line.split(None, 1)[0])
                if Test:
                    print "would have killed: ", app
                else:
                    print "killing ", app
                    os.system("sudo kill -%d %s" % (signal.SIGKILL, pid))
            except:
                print "Unable to kill process ", app 

def remove_files(files):
    for f in files:
        if '~' in f:
            f.replace('~', os.path.expanduser('~'))
            if Test:
                print "would have removed: ", f
            else:
                print "removing: ", f
                os.remove(f)
        else:
            if Test:
                print "would have removed: ", f
            else:
                print "removing: :, f
                os.system("sudo rm -rf %s" % (f,))
    

def do_one(malware):
    try:
        for i in WARE_MAPPING[malware]['order']:
            if i == 'files':
                remove_files(WARE_MAPPING[malware]['files'])
            if i == 'kill':
                for k in WARE_MAPPING[malware]['kill']:
                    kill_app(i)
    except:
        print "Unable to find malware to remove"

def do_all():
    for j in WARE_MAPPING.keys():
        do_one(j)


if __name__ == '__main__':

    try:
        args = sys.argv[1]
    except:
        print "please specify ", ','.join(WARE_MAPPING.keys()), "or all"
        print "you can also use test as the second argument to see what would happen"
        sys.exit(1)
    try:
        if sys.argv[2]:
            Test = True
    except:
        pass
    print "Scanning System for Malicious Files"
    if args == 'all':
        do_all()
    else:
        do_one(args)
    
