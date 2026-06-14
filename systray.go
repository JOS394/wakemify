//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

static id g_statusItem = nil;

extern void goOnToggleWindow(void);
extern void setWindowPosition(double, double);

@interface WakemifyHandler : NSObject
@end

@implementation WakemifyHandler
- (void)statusItemAction {
	NSWindow *btnWindow = [[g_statusItem button] window];
	double x = 0, y = 0;
	if (btnWindow) {
		NSRect btnFrame = [[g_statusItem button] frame];
		NSRect windowFrame = [btnWindow frame];
		x = windowFrame.origin.x + btnFrame.origin.x;
		y = windowFrame.origin.y + windowFrame.size.height - 2;
	} else {
		NSScreen *screen = [NSScreen mainScreen];
		x = [screen frame].size.width - 240;
		y = [screen frame].size.height - 2;
	}
	setWindowPosition(x, y);
	goOnToggleWindow();
}
@end

void setupStatusBar(void) {
	void (^block)(void) = ^{
		@autoreleasepool {
			g_statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
			id handler = [[WakemifyHandler alloc] init];
			[g_statusItem setAction:@selector(statusItemAction)];
			[g_statusItem setTarget:handler];
		}
	};
	if ([NSThread isMainThread]) { block(); }
	else { dispatch_sync(dispatch_get_main_queue(), block); }
}

static void do_setStatusIcon(const unsigned char* data, int len) {
	NSData *nsData = [NSData dataWithBytes:data length:len];
	NSImage *image = [[NSImage alloc] initWithData:nsData];
	[image setTemplate:YES];
	[g_statusItem setImage:image];
	[image release];
}
void setStatusIcon(const unsigned char* data, int len) {
	if ([NSThread isMainThread]) { do_setStatusIcon(data, len); }
	else {
		unsigned char* d = malloc(len); memcpy(d, data, len);
		dispatch_sync(dispatch_get_main_queue(), ^{ do_setStatusIcon(d, len); free(d); });
	}
}
*/
import "C"
import "embed"

//go:embed icons/iconTemplate.png icons/iconActiveTemplate.png
var iconFS embed.FS
var smInstance *SystrayManager

type SystrayManager struct {
	app *App
}

func NewSystrayManager(app *App) *SystrayManager {
	return &SystrayManager{app: app}
}

func (sm *SystrayManager) Run() {
	smInstance = sm
	C.setupStatusBar()
	sm.UpdateUI(false)
}

func (sm *SystrayManager) UpdateUI(active bool) {
	icon := "icons/iconTemplate.png"
	if active {
		icon = "icons/iconActiveTemplate.png"
	}
	data, _ := iconFS.ReadFile(icon)
	C.setStatusIcon((*C.uchar)(&data[0]), C.int(len(data)))
}
