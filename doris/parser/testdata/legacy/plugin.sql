INSTALL PLUGIN FROM "/home/users/doris/auditdemo.zip";
INSTALL PLUGIN FROM "/home/users/doris/auditdemo/";
INSTALL PLUGIN FROM "http://mywebsite.com/plugin.zip";
INSTALL PLUGIN FROM "http://mywebsite.com/plugin.zip" PROPERTIES("md5sum" = "73877f6029216f4314d712086a146570");
UNINSTALL PLUGIN auditdemo;
SHOW PLUGINS
