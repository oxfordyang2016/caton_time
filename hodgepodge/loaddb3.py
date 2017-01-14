#loaddb3.mysql_load2('/loadtest2/','localhost',3306,'root','yangmingtestmysql')
import os
dbstring1='--database cydex'
dbstring2='--database cydex_user'
def mysql_load2(dir,host,port,user,passwd):
    pure_file_list=filter(lambda x:(os.path.isfile(x)and('bin' in x or 'sql' in x)) ,map(lambda x:os.path.join(dir,x),os.listdir(dir)))
    time_file_dict={}
    file_create_time_list=map(lambda x:file_time_dict[os.path.getctime(x)],pure_file_list)
    time_file_dict = dict(zip(file_create_time_list,pure_file_list ))
    backup_dir_path=[time_file_dict[key] for key in sorted(time_file_dict.iterkeys())]
    def back(path):
    	try:
	        if "bin" in path:
	            commandline='mysqlbinlog '+path+'| mysql -h '+str(host)+' --port '+str(port)+' -u '+str(user)+' -p'+str(passwd)
	        elif "sql" in key:
	            commandline='mysql -h '+str(host)+' --port '+str(port)+ ' -u '+str(user)+' -p'+ str(passwd)+ ' < '+path
	        os.system(commandline)
	        if os.path.exists(os.path.join(dir,'backup'))==False:
	        	os.mkdir(os.path.join(dir,'backup'))
	        os.system('mv '+path+' '+os.path.join(dir,'backup'))	
            return (str(path)+'ok')
         except:
            return(str(path)+'error') 
    back_log=map(back,backup_dir_path) 