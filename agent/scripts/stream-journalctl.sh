# CRONJOB -  */10 * * * * /var/log/stream-journalctl.sh 
journalctl --since "10 minutes ago" >> /var/log/magpie-monitor/journal.log