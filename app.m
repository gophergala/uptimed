#include <Cocoa/Cocoa.h>

NSStatusItem *statusItem;

int StartApp(void) {
  [NSAutoreleasePool new];
  [NSApplication sharedApplication];
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

  NSMenuItem *tItem = nil;
  NSMenu *theMenu;

  theMenu = [[NSMenu alloc] initWithTitle:@""];
  [theMenu setAutoenablesItems:NO];
  tItem = [theMenu addItemWithTitle:@"Quit" action:@selector(terminate:) keyEquivalent:@"q"];
  [tItem setKeyEquivalentModifierMask:NSCommandKeyMask];

  NSStatusBar *statusBar = [NSStatusBar systemStatusBar];
  statusItem = [statusBar statusItemWithLength:NSVariableStatusItemLength];
  [statusItem retain];
  [statusItem setTitle:@"uptimed"];
  [statusItem setHighlightMode:YES];
  [statusItem setMenu:theMenu];

  [NSApp activateIgnoringOtherApps:YES];
  [NSApp run];
  return 0;
}

void SetLabelText(const char *str) {
  @autoreleasepool {
    NSString *text = [NSString stringWithUTF8String:str];
    [statusItem setTitle:text];
  }
}
