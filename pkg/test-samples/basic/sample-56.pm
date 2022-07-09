try
	print(variable)
except UndefinedIdentifier, AnyException as errors
	print(errors)
except NoToStringException as errors
	print(errors)
else
	print("Unknown *errors")
	raise UnknownException()
finally
	print("Done")
end