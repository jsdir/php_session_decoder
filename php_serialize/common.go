package php_serialize

const (
	TOKEN_NULL              rune = 'N'
	TOKEN_BOOL              rune = 'b'
	TOKEN_INT               rune = 'i'
	TOKEN_FLOAT             rune = 'd'
	TOKEN_STRING            rune = 's'
	TOKEN_ARRAY             rune = 'a'
	TOKEN_OBJECT            rune = 'O'
	TOKEN_OBJECT_SERIALIZED rune = 'C'
	TOKEN_REFERENCE         rune = 'R'
	TOKEN_REFERENCE_OBJECT  rune = 'r'
	TOKEN_SPL_ARRAY         rune = 'x'
	TOKEN_SPL_ARRAY_MEMBERS rune = 'm'

	SEPARATOR_VALUE_TYPE rune = ':'
	SEPARATOR_VALUES     rune = ';'

	DELIMITER_STRING_LEFT  rune = '"'
	DELIMITER_STRING_RIGHT rune = '"'
	DELIMITER_OBJECT_LEFT  rune = '{'
	DELIMITER_OBJECT_RIGHT rune = '}'

	FORMATTER_FLOAT     byte = 'g'
	FORMATTER_PRECISION int  = 17
)

var (
	debugMode = false
)

func Debug(value bool) {
	debugMode = value
}

func NewPhpObject(className string) *PhpObject {
	return &PhpObject{
		className: className,
		members:   NewPhpArray(),
	}
}

type SerializedDecodeFunc func(string) (PhpValue, error)

type SerializedEncodeFunc func(PhpValue) (string, error)

type PhpValue interface{}

type PhpArray struct {
	keys   []PhpValue
	values map[PhpValue]PhpValue
}

func NewPhpArrayFromData(data map[interface{}]interface{}) *PhpArray {
	phpArray := NewPhpArray()

	for k, v := range data {
		phpArray.keys = append(phpArray.keys, k)
		phpArray.values[k] = v
	}

	return phpArray
}

func NewPhpArray() *PhpArray {
	return &PhpArray{values: map[PhpValue]PhpValue{}}
}

func (a *PhpArray) Set(key PhpValue, value PhpValue) {
	if _, ok := a.values[key]; !ok {
		a.keys = append(a.keys, key)
	}

	a.values[key] = value
}

func (a *PhpArray) Get(key PhpValue) (value PhpValue, ok bool) {
	value, ok = a.values[key]
	return
}

func (a *PhpArray) Keys() (keys []PhpValue) {
	for _, k := range a.keys {
		keys = append(keys, k)
	}

	return
}

func (a *PhpArray) Delete(key PhpValue) {
	newKeys := []PhpValue{}
	for _, k := range a.keys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}

	a.keys = newKeys
	delete(a.values, key)
}

func (a *PhpArray) ReplaceKey(oldKey PhpValue, newKey PhpValue) {
	for i, k := range a.keys {
		if k == oldKey {
			a.keys[i] = newKey
			a.values[newKey] = a.values[oldKey]
			delete(a.values, oldKey)
			return
		}
	}
}

type PhpSlice []PhpValue

type PhpObject struct {
	className string
	members   *PhpArray
}

func (self *PhpObject) GetClassName() string {
	return self.className
}

func (self *PhpObject) SetClassName(name string) *PhpObject {
	self.className = name
	return self
}

func (self *PhpObject) GetMembers() *PhpArray {
	return self.members
}

func (self *PhpObject) SetMembers(members *PhpArray) *PhpObject {
	self.members = members
	return self
}

func (self *PhpObject) GetPrivate(name string) (v PhpValue, ok bool) {
	return self.members.Get("\x00" + self.className + "\x00" + name)
}

func (self *PhpObject) SetPrivate(name string, value PhpValue) *PhpObject {
	self.members.Set("\x00"+self.className+"\x00"+name, value)
	return self
}

func (self *PhpObject) GetProtected(name string) (v PhpValue, ok bool) {
	return self.members.Get("\x00*\x00" + name)
}

func (self *PhpObject) SetProtected(name string, value PhpValue) *PhpObject {
	self.members.Set("\x00*\x00"+name, value)
	return self
}

func (self *PhpObject) GetPublic(name string) (v PhpValue, ok bool) {
	return self.members.Get(name)
}

func (self *PhpObject) SetPublic(name string, value PhpValue) *PhpObject {
	self.members.Set(name, value)
	return self
}

func NewPhpObjectSerialized(className string) *PhpObjectSerialized {
	return &PhpObjectSerialized{
		className: className,
	}
}

type PhpObjectSerialized struct {
	className string
	data      string
	value     PhpValue
}

func (self *PhpObjectSerialized) GetClassName() string {
	return self.className
}

func (self *PhpObjectSerialized) SetClassName(name string) *PhpObjectSerialized {
	self.className = name
	return self
}

func (self *PhpObjectSerialized) GetData() string {
	return self.data
}

func (self *PhpObjectSerialized) SetData(data string) *PhpObjectSerialized {
	self.data = data
	return self
}

func (self *PhpObjectSerialized) GetValue() PhpValue {
	return self.value
}

func (self *PhpObjectSerialized) SetValue(value PhpValue) *PhpObjectSerialized {
	self.value = value
	return self
}

func NewPhpSplArray(array, properties PhpValue) *PhpSplArray {
	if array == nil {
		array = NewPhpArray()
	}

	if properties == nil {
		properties = NewPhpArray()
	}

	return &PhpSplArray{
		array:      array,
		properties: properties,
	}
}

type PhpSplArray struct {
	flags      int
	array      PhpValue
	properties PhpValue
}

func (self *PhpSplArray) GetFlags() int {
	return self.flags
}

func (self *PhpSplArray) SetFlags(value int) {
	self.flags = value
}

func (self *PhpSplArray) GetArray() PhpValue {
	return self.array
}

func (self *PhpSplArray) SetArray(value PhpValue) {
	self.array = value
}

func (self *PhpSplArray) GetProperties() PhpValue {
	return self.properties
}

func (self *PhpSplArray) SetProperties(value PhpValue) {
	self.properties = value
}
