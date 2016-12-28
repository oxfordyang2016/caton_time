def pre(q):
    if q>=5:
        return 'true'


print [k*k for k in range(10) if pre(k) ]        