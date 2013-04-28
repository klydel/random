### Random

Collection of Random stuff

snapshot_hosts.py
-----------

Script to perform backups on amazon ec2 hosts using snapshots.  Add hosts to SNAPSHOT_HOSTS if host needs to be backed up.  if a host has multiple volumes, a number will be postfixed to the snapshot name.  This script can also report backup status to nagios using NAGIOS_CMD_FILE.  if nagios isnt required, change NAGIOS_CMD_FILE to a log file.  
<pre><code>
    snapshot_hosts.py

</code></pre>






