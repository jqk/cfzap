package cfzap

// ConfigOption defines the information needed to create a logger from a config file.
type ConfigOption struct {
	// CreateNew indicates if a new logger should be created.
	CreateNew bool
	// FileName is the config file name without extension, default is 'cfzap'.
	FileName string
	// FileExt is the config file extension, default is empty string.
	// The supported extensions are:
	// "json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "dotenv", "env" and "ini".
	// If this field is set to empty string (default), the extensions be searched one by one according to above list.
	// It's better giving a non empty value to this field if 'FileName' has extension.
	FileExt string
	// FilePaths is the path list that the config file may be located.
	// If this field is nil or empty, current executable path, current path will be used.
	FilePaths []string
}

// ConfigPropertySetter defines the function type to change ConfigOption's property.
type ConfigPropertySetter func(*ConfigOption)

// WithCreateNew set up the CreateNew property of a ConfigOption object。
func WithCreateNew(createNew bool) ConfigPropertySetter {
	return func(option *ConfigOption) {
		option.CreateNew = createNew
	}
}

// WithFileName set up the FileName property of a ConfigOption object。
func WithFileName(fileName string) ConfigPropertySetter {
	return func(option *ConfigOption) {
		option.FileName = fileName
	}
}

// WithFileExt set up the FileExt property of a ConfigOption object。
func WithFileExt(fileExt string) ConfigPropertySetter {
	return func(option *ConfigOption) {
		option.FileExt = fileExt
	}
}

// WithFilePaths set up the FilePaths property of a ConfigOption object。
func WithFilePaths(filePaths ...string) ConfigPropertySetter {
	return func(option *ConfigOption) {
		option.FilePaths = filePaths
	}
}

// NewConfigOption creates and returns ConfigOption object.
func NewConfigOption(setters ...ConfigPropertySetter) *ConfigOption {
	// create default value.
	option := &ConfigOption{
		CreateNew: false,
		FileName:  "cfzap",
		FileExt:   "",
		FilePaths: []string{"."}}

	// apply settings to properties.
	for _, setter := range setters {
		setter(option)
	}

	return option
}
