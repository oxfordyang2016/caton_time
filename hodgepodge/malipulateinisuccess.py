#i am using python 2.6
#http://stackoverflow.com/questions/8884188/how-to-read-and-write-ini-file-with-python
try:
    from configparser import ConfigParser
except ImportError:
    from ConfigParser import ConfigParser  # ver. < 3.0

# instantiate
config = ConfigParser()

# parse existing file
config.read('/tmp/goodday/test.ini')

# read values from a section
string_val = config.get('section_a', 'string_val')
bool_val = config.getboolean('section_a', 'bool_val')
int_val = config.getint('section_a', 'int_val')
float_val = config.getfloat('section_a', 'pi_val')

# update existing value
config.set('section_a', 'string_val', 'worldisfutureawesomeday')

# add a new section and some values
config.add_section('sectionyangming')
config.set('sectionyangming', 'meal_val', 'spam')
config.set('sectionyangming', 'not_found_val','yangming')

# save to a file
with open('/tmp/goodday/test.ini', 'w') as configfile:#you can modify the position
    config.write(configfile)