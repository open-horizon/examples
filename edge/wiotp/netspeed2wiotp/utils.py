## utils.py
# Useful functions for logging: print_
import os
import sys

# print_ and other important builtins
try:
    import builtins
except ImportError:
    def print_(*args, **kwargs):
        """The new-style print function taken from https://pypi.python.org/pypi/six/ """
        fp = kwargs.pop("file", sys.stdout)
        if fp is None:
            return

        def write(data):
            if not isinstance(data, basestring):
                data = str(data)
            fp.write(data)

        want_unicode = False
        sep = kwargs.pop("sep", None)
        if sep is not None:
            if isinstance(sep, unicode):
                want_unicode = True
            elif not isinstance(sep, str):
                raise TypeError("sep must be None or a string")
        end = kwargs.pop("end", None)
        if end is not None:
            if isinstance(end, unicode):
                want_unicode = True
            elif not isinstance(end, str):
                raise TypeError("end must be None or a string")
        if kwargs:
            raise TypeError("invalid keyword arguments to print()")
        if not want_unicode:
            for arg in args:
                if isinstance(arg, unicode):
                    want_unicode = True
                    break
        if want_unicode:
            newline = unicode("\n")
            space = unicode(" ")
        else:
            newline = "\n"
            space = " "
        if sep is None:
            sep = space
        if end is None:
            end = newline
        for i, arg in enumerate(args):
            if i:
                write(sep)
            write(arg)
        write(end)
else:
    print_ = getattr(builtins, 'print')
    del builtins

def get_serial(debug_flag=False):
    # Extract serial from cpuinfo file
    cpuserial = "0000000000000000"
    try:
        f = open('/proc/cpuinfo','r')
        for line in f:
            if line[0:6]=='Serial':
               cpuserial = line[10:26]
        f.close()
    except:
        cpuserial = "ERROR000000000"
    
    if debug_flag:
        print_('processor s/n: %s' % cpuserial)

    return cpuserial

def check_env_var(envname, default='', printerr=True):
    ''' '''
    if envname in os.environ:
       val = os.getenv(envname)
       if val == '' or val == '-':
           if printerr:
               print_("utils.py: variable" + envname + " value is '%s'" % val)
           return default
       return val
    
    else:
       if printerr:
           print_("utils.py: Environment variable " + envname + " not found.")
       return default

def getEnvInt(name, default):
    """Return the named env var value as an int, or the default value."""
    if name in os.environ:
        strVal = os.getenv(name)
        try:
            return int(strVal)
        except ValueError as e:
            print_('workload_config.py: Error: invalid value for environment variable %s: %s. Using default value %d.' % (name, str(e), default) )
            return default
    else:
        return default