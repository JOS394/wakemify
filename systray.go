//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

static id g_statusItem = nil;
static id g_handler = nil;

extern void goOnToggleWindow(void);
extern void setWindowPosition(double, double);

@interface WakemifyHandler : NSObject
@end

@implementation WakemifyHandler
- (void)statusItemAction {
	double x = 0, y = 0;
	NSScreen *screen = [NSScreen mainScreen];
	if (!screen) screen = [[NSScreen screens] firstObject];
	NSRect visibleFrame = [screen visibleFrame];
	id button = [g_statusItem button];
	NSWindow *btnWindow = [button window];
	if (button && btnWindow) {
		NSRect btnFrame = [button frame];
		NSRect winFrame = [btnWindow frame];
		x = (winFrame.origin.x + btnFrame.origin.x) - visibleFrame.origin.x;
		y = 2;
	} else {
		x = visibleFrame.size.width - 240;
		y = 2;
	}
	setWindowPosition(x, y);
	goOnToggleWindow();
}
@end

void setupStatusBar(void) {
	void (^block)(void) = ^{
		@autoreleasepool {
			g_statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
			g_handler = [[WakemifyHandler alloc] init];
			[g_statusItem setAction:@selector(statusItemAction)];
			[g_statusItem setTarget:g_handler];
			[[g_statusItem button] sendActionOn:NSEventMaskLeftMouseDown];
		}
	};
	if ([NSThread isMainThread]) { block(); }
	else { dispatch_sync(dispatch_get_main_queue(), block); }
}

static void do_setStatusIcon(const unsigned char* data, int len) {
	NSData *nsData = [NSData dataWithBytes:data length:len];
	NSImage *image = [[NSImage alloc] initWithData:nsData];
	[image setTemplate:YES];
	[image setSize:NSMakeSize(22, 22)];
	[g_statusItem setImage:image];
	[image release];
}

static void do_removeButtons(void) {
	NSArray *windows = [NSApp windows];
	for (NSWindow *win in windows) {
		NSString *className = NSStringFromClass([win class]);
		if ([className containsString:@"StatusBar"] || [className containsString:@"NSStatus"]) continue;
		[[win standardWindowButton:NSWindowCloseButton] setHidden:YES];
		[[win standardWindowButton:NSWindowMiniaturizeButton] setHidden:YES];
		[[win standardWindowButton:NSWindowZoomButton] setHidden:YES];
	}
}
void removeStandardWindowButtons(void) {
	dispatch_async(dispatch_get_main_queue(), ^{ do_removeButtons(); });
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

func removeWindowButtons() {
	C.removeStandardWindowButtons()
}

func (sm *SystrayManager) UpdateUI(active bool) {
	icon := "icons/iconTemplate.png"
	if active {
		icon = "icons/iconActiveTemplate.png"
	}
	data, _ := iconFS.ReadFile(icon)
	C.setStatusIcon((*C.uchar)(&data[0]), C.int(len(data)))
}
