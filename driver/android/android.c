// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on https://github.com/golang/mobile
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go:build android
//  +build android

#include <android/log.h>
#include <dlfcn.h>
#include <errno.h>
#include <fcntl.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include "_cgo_export.h"

#define LOG_INFO(...) __android_log_print(ANDROID_LOG_INFO, "Go", __VA_ARGS__)
#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Go", __VA_ARGS__)

static jclass current_class;

static jclass find_class(JNIEnv *env, const char *class_name)
{
	jclass clazz = (*env)->FindClass(env, class_name);
	if (clazz == NULL)
	{
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find %s", class_name);
		return NULL;
	}
	return clazz;
}

static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
	jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
	if (m == 0)
	{
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

static jmethodID find_static_method(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
	jmethodID m = (*env)->GetStaticMethodID(env, clazz, name, sig);
	if (m == 0)
	{
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

static jmethodID key_rune_method;
static jmethodID show_keyboard_method;
static jmethodID hide_keyboard_method;

jint JNI_OnLoad(JavaVM *vm, void *reserved)
{
	JNIEnv *env;
	if ((*vm)->GetEnv(vm, (void **)&env, JNI_VERSION_1_6) != JNI_OK)
	{
		return -1;
	}

	return JNI_VERSION_1_6;
}

static int main_running = 0;

// Entry point from our subclassed NativeActivity.
//
// By here, the Go runtime has been initialized (as we are running in
// -buildmode=c-shared) but the first time it is called, Go's main.main
// hasn't been called yet.
//
// The Activity may be created and destroyed multiple times throughout
// the life of a single process. Each time, onCreate is called.
void ANativeActivity_onCreate(ANativeActivity *activity, void *savedState, size_t savedStateSize)
{
	if (!main_running)
	{
		JNIEnv *env = activity->env;

		// Note that activity->clazz is mis-named.
		current_class = (*env)->GetObjectClass(env, activity->clazz);
		current_class = (*env)->NewGlobalRef(env, current_class);
		key_rune_method = find_static_method(env, current_class, "getRune", "(III)I");
		show_keyboard_method = find_static_method(env, current_class, "showKeyboard", "(I)V");
		hide_keyboard_method = find_static_method(env, current_class, "hideKeyboard", "()V");

		setCurrentContext(activity->vm, (*env)->NewGlobalRef(env, activity->clazz));

		// Set FILESDIR
		if (setenv("FILESDIR", activity->internalDataPath, 1) != 0)
		{
			LOG_INFO("setenv(\"FILESDIR\", \"%s\", 1) failed: %d", activity->internalDataPath, errno);
		}

		// Set TMPDIR.
		jmethodID gettmpdir = find_method(env, current_class, "getTmpdir", "()Ljava/lang/String;");
		jstring jpath = (jstring)(*env)->CallObjectMethod(env, activity->clazz, gettmpdir, NULL);
		const char *tmpdir = (*env)->GetStringUTFChars(env, jpath, NULL);
		if (setenv("TMPDIR", tmpdir, 1) != 0)
		{
			LOG_INFO("setenv(\"TMPDIR\", \"%s\", 1) failed: %d", tmpdir, errno);
		}
		(*env)->ReleaseStringUTFChars(env, jpath, tmpdir);

		// Call the Go main.main.
		uintptr_t mainPC = (uintptr_t)dlsym(RTLD_DEFAULT, "main.main");
		if (!mainPC)
		{
			LOG_FATAL("missing main.main");
		}
		callMain(mainPC);
		main_running = 1;
	}

	// These functions match the methods on Activity, described at
	// http://developer.android.com/reference/android/app/Activity.html
	//
	// Note that onNativeWindowResized is not called on resize. Avoid it.
	// https://code.google.com/p/android/issues/detail?id=180645
	activity->callbacks->onStart = onStart;
	activity->callbacks->onResume = onResume;
	activity->callbacks->onSaveInstanceState = onSaveInstanceState;
	activity->callbacks->onPause = onPause;
	activity->callbacks->onStop = onStop;
	activity->callbacks->onDestroy = onDestroy;
	activity->callbacks->onWindowFocusChanged = onWindowFocusChanged;
	activity->callbacks->onNativeWindowCreated = onNativeWindowCreated;
	activity->callbacks->onNativeWindowRedrawNeeded = onNativeWindowRedrawNeeded;
	activity->callbacks->onNativeWindowDestroyed = onNativeWindowDestroyed;
	activity->callbacks->onInputQueueCreated = onInputQueueCreated;
	activity->callbacks->onInputQueueDestroyed = onInputQueueDestroyed;
	activity->callbacks->onConfigurationChanged = onConfigurationChanged;
	activity->callbacks->onLowMemory = onLowMemory;

	onCreate(activity);
}

int32_t getKeyRune(JNIEnv *env, AInputEvent *e)
{
	return (int32_t)(*env)->CallStaticIntMethod(
		env,
		current_class,
		key_rune_method,
		AInputEvent_getDeviceId(e),
		AKeyEvent_getKeyCode(e),
		AKeyEvent_getMetaState(e));
}

void showKeyboard(JNIEnv *env, int keyboardType)
{
	(*env)->CallStaticVoidMethod(
		env,
		current_class,
		show_keyboard_method,
		keyboardType);
}

void hideKeyboard(JNIEnv *env)
{
	(*env)->CallStaticVoidMethod(
		env,
		current_class,
		hide_keyboard_method);
}

void Java_org_golang_app_GoNativeActivity_insetsChanged(JNIEnv *env, jclass clazz, int top, int bottom, int left, int right)
{
	insetsChanged(top, bottom, left, right);
}

void Java_org_golang_app_GoNativeActivity_keyboardTyped(JNIEnv *env, jclass clazz, jstring str)
{
	const char *cstr = (*env)->GetStringUTFChars(env, str, JNI_FALSE);
	keyboardTyped((char *)cstr);
}

void Java_org_golang_app_GoNativeActivity_keyboardDelete(JNIEnv *env, jclass clazz)
{
	keyboardDelete();
}

void Java_org_golang_app_GoNativeActivity_setDarkMode(JNIEnv *env, jclass clazz, jboolean dark)
{
	setDarkMode((bool)dark);
}

void Java_org_golang_app_GoNativeActivity_scrolled(JNIEnv *env, jclass clazz, jfloat posX, jfloat posY, jfloat distanceX, jfloat distanceY)
{
	scrolled((float)posX, (float)posY, (float)distanceX, (float)distanceY);
}

void Java_org_golang_app_GoNativeActivity_scaled(JNIEnv *env, jclass clazz, jfloat scaleFactor, jfloat posX, jfloat posY)
{
	scaled((float)scaleFactor, (float)posX, (float)posY);
}

void Java_org_golang_app_GoNativeActivity_longPressed(JNIEnv *env, jclass clazz, jfloat posX, jfloat posY)
{
	longPressed((float)posX, (float)posY);
}
