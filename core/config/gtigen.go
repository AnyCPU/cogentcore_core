// Code generated by "core generate"; DO NOT EDIT.

package config

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Config", IDName: "config", Doc: "Config is the main config struct\nthat contains all of the configuration\noptions for the Cogent Core tool", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "Name", Doc: "the user-friendly name of the project"}, {Name: "ID", Doc: "the bundle / package ID to use of the project (required for building for mobile platforms\nand packaging for desktop platforms). It is typically in the format com.org.app (eg: com.core.mail)"}, {Name: "Desc", Doc: "the description of the project"}, {Name: "Version", Doc: "the version of the project"}, {Name: "Type", Doc: "the type of the project (app/library)"}, {Name: "Build", Doc: "the configuration options for the build, install, run, and pack commands"}, {Name: "Pack", Doc: "the configuration information for the pack command"}, {Name: "Web", Doc: "the configuration information for web"}, {Name: "Setup", Doc: "the configuration options for the setup command"}, {Name: "Log", Doc: "the configuration options for the log command"}, {Name: "Release", Doc: "the configuration options for the release command"}, {Name: "Generate", Doc: "the configuration options for the generate command"}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Build", IDName: "build", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "Package", Doc: "the path of the package to build"}, {Name: "Target", Doc: "the target platforms to build executables for"}, {Name: "Output", Doc: "the output file name; if not specified, it depends on the package being built"}, {Name: "Debug", Doc: "whether to build/run the app in debug mode; this currently only works on mobile platforms"}, {Name: "Rebuild", Doc: "force rebuilding of packages that are already up-to-date"}, {Name: "Install", Doc: "install the generated executable"}, {Name: "PrintOnly", Doc: "print the commands but do not run them"}, {Name: "Print", Doc: "print the commands"}, {Name: "GCFlags", Doc: "arguments to pass on each go tool compile invocation"}, {Name: "LDFlags", Doc: "arguments to pass on each go tool link invocation"}, {Name: "Tags", Doc: "a comma-separated list of additional build tags to consider satisfied during the build"}, {Name: "Trimpath", Doc: "remove all file system paths from the resulting executable. Instead of absolute file system paths, the recorded file names will begin either a module path@version (when using modules), or a plain import path (when using the standard library, or GOPATH)."}, {Name: "Work", Doc: "print the name of the temporary work directory and do not delete it when exiting"}, {Name: "IOSVersion", Doc: "the minimal version of the iOS SDK to compile against"}, {Name: "AndroidMinSDK", Doc: "the minimum supported Android SDK (uses-sdk/android:minSdkVersion in AndroidManifest.xml)"}, {Name: "AndroidTargetSDK", Doc: "the target Android SDK version (uses-sdk/android:targetSdkVersion in AndroidManifest.xml)"}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Pack", IDName: "pack", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "DMG", Doc: "whether to build a .dmg file on macOS in addition to a .app file.\nThis is automatically disabled for the install command."}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Setup", IDName: "setup", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "Platform", Doc: "the platform to set things up for"}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Log", IDName: "log", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "Target", Doc: "the target platform to view the logs for (ios or android)"}, {Name: "Keep", Doc: "whether to keep the previous log messages or clear them"}, {Name: "All", Doc: "messages not generated from your app equal to or above this log level will be shown"}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Release", IDName: "release", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "VersionFile", Doc: "the Go file to store version information in"}, {Name: "Package", Doc: "the Go package in which the version file will be stored"}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Generate", IDName: "generate", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "Enumgen", Doc: "the enum generation configuration options passed to enumgen"}, {Name: "Gtigen", Doc: "the generation configuration options passed to gtigen"}, {Name: "Dir", Doc: "the source directory to run generate on (can be multiple through ./...)"}, {Name: "Output", Doc: "the output file location relative to the package on which generate is being called"}}})

var _ = gti.AddType(&gti.Type{Name: "cogentcore.org/core/core/config.Web", IDName: "web", Doc: "Web containts the configuration information for building for web and creating\nthe HTML page that loads a Go wasm app and its resources.", Directives: []gti.Directive{{Tool: "gti", Directive: "add"}}, Fields: []gti.Field{{Name: "Port", Doc: "Port is the port to serve the page at when using the serve command."}, {Name: "RandomVersion", Doc: "RandomVersion is whether to automatically add a random string to the\nend of the version string for the app when building for web. This is\nnecessary in order for changes made during local development to show up,\nbut should not be enabled in release builds to prevent constant inaccurate\nupdate messages. It is enabled by default in the serve command and disabled\nby default otherwise."}, {Name: "Gzip", Doc: "Gzip is whether to gzip the app.wasm file that is built in the build command\nand serve it as a gzip-encoded file in the run command."}, {Name: "BackgroundColor", Doc: "A placeholder background color for the application page to display before\nits stylesheets are loaded.\n\nDEFAULT: #2d2c2c."}, {Name: "ThemeColor", Doc: "The theme color for the application. This affects how the OS displays the\napp (e.g., PWA title bar or Android's task switcher).\n\nDEFAULT: #2d2c2c."}, {Name: "LoadingLabel", Doc: "The text displayed while loading a page. Load progress can be inserted by\nincluding \"{progress}\" in the loading label.\n\nDEFAULT: \"{progress}%\"."}, {Name: "Lang", Doc: "The page language.\n\nDEFAULT: en."}, {Name: "Author", Doc: "The page authors."}, {Name: "Keywords", Doc: "The page keywords."}, {Name: "Image", Doc: "The path of the default image that is used by social networks when\nlinking the app."}, {Name: "AutoUpdateInterval", Doc: "The interval between each app auto-update while running in a web browser.\nZero or negative values deactivates the auto-update mechanism.\n\nDefault is 10 seconds."}, {Name: "Env", Doc: "The environment variables that are passed to the progressive web app.\n\nReserved keys:\n- GOAPP_VERSION\n- GOAPP_GOAPP_STATIC_RESOURCES_URL"}, {Name: "WasmContentLengthHeader", Doc: "The HTTP header to retrieve the WebAssembly file content length.\n\nContent length finding falls back to the Content-Length HTTP header when\nno content length is found with the defined header."}}})
