package heap

import (
	"jvmgo/ch11/classpath"
	"fmt"
	"jvmgo/ch11/classfile"
)

type ClassLoader struct {
	cp          *classpath.Classpath
	verboseFlag bool
	classMap    map[string]*Class //loaded classes
}

func NewClassLoader(cp *classpath.Classpath, verboseFlag bool) *ClassLoader {
	loader := &ClassLoader{
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
	}
	loader.loadBasicClasses()
	loader.loadPrimitiveClasses()
	return loader
}


func (classLoader *ClassLoader) loadBasicClasses(){
	jlClassClass := classLoader.LoadClass("java/lang/Class")
	for _, class := range classLoader.classMap {
		if class.jClass == nil {
			class.jClass = jlClassClass.NewObject()
			class.jClass.extra = class
		}
	}
}

func (classLoader *ClassLoader) loadPrimitiveClasses(){
	for primitiveType, _ := range primitiveTypes {
		classLoader.loadPrimitiveClass(primitiveType)
	}
}

func (classLoader *ClassLoader) loadPrimitiveClass(className string)  {
	class := &Class{
		accessFlags:ACC_PUBLIC,
		name : className,
		loader:classLoader,
		initStarted:true,
	}
	class.jClass = classLoader.classMap["java/lang/Class"].NewObject()
	class.jClass.extra = class
	classLoader.classMap[className] = class
}

func (classLoader *ClassLoader) LoadClass(name string) *Class {
	if class, ok := classLoader.classMap[name]; ok {
		return class
	}
	var class *Class
	if name[0] == '[' {
		class = classLoader.loadArrayClass(name)
	}else {
		class = classLoader.loadNonArrayClass(name)
	}
	if jlClassClass, ok := classLoader.classMap["java/lang/Class"]; ok{
		class.jClass = jlClassClass.NewObject()
		class.jClass.extra = class
	}
	return class
}

func (classLoader *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry := classLoader.readClass(name)
	class := classLoader.defineClass(data)
	link(class)

	if classLoader.verboseFlag {
		fmt.Printf("[Loaded %s from %s]\n", name, entry)
	}

	return class
}

func (classLoader *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := classLoader.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}

func (classLoader *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	hackClass(class)
	class.loader = classLoader
	resolveSuperClass(class)
	resolveInterfaces(class)
	classLoader.classMap[class.name] = class
	return class
}

func (classLoader *ClassLoader) loadArrayClass(name string) *Class {
	class := &Class{
		accessFlags: ACC_PUBLIC,
		name:        name,
		loader:      classLoader,
		initStarted: true,
		superClass:  classLoader.LoadClass("java/lang/Object"),
		interfaces: []*Class{
			classLoader.LoadClass("java/lang/Cloneable"),
			classLoader.LoadClass("java/io/Serializable"),
		},
	}
	classLoader.classMap[name] = class
	return class

}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		panic("java.lang.ClassFormatError")
	}
	return newClass(cf)
}

func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}

func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

func link(class *Class) {
	verify(class)
	prepare(class)
}

func verify(class *Class) {

}

func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}

func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}

func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}

func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initStaticFinalVar(class, field)
		}
	}
}

func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex()
	slotId := field.SlotId()

	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			goStr := cp.GetConstant(cpIndex).(string)
			jStr := JString(class.Loader(), goStr)
			vars.SetRef(slotId,jStr)
		}
	}
}

// todo
func hackClass(class *Class) {
	if class.name == "java/lang/ClassLoader" {
		loadLibrary := class.GetStaticMethod("loadLibrary", "(Ljava/lang/Class;Ljava/lang/String;Z)V")
		loadLibrary.code = []byte{0xb1} // return void
	}
}