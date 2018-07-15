
def hello(fn):
	def wrapper():
		print "hello, %s" % fn.__name__
		fn()
		print "goodbye, %s" % fn.__name__

	return wrapper

@hello
def Hao():
	print "I am Hao Chen"

Hao()