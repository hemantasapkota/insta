package ios

const (
	//StateDidFinishLaunchingWithOptions ...
	StateDidFinishLaunchingWithOptions string = "didFinishLaunchingWithOptions"
	//StateDidRegisterForRemoteNotificationsWithDeviceToken ...
	StateDidRegisterForRemoteNotificationsWithDeviceToken string = "didRegisterForRemoteNotificationsWithDeviceToken"
	//StateDidBecomeActive ...
	StateDidBecomeActive string = "didBecomeActive"
	//StateDidEnterBackground ...
	StateDidEnterBackground string = "didEnterBackground"
	//StateWillEnterForeground ...
	StateWillEnterForeground string = "willEnterBackground"
	//StateWillTerminate ...
	StateWillTerminate string = "willTerminate"
	//StateAppOpenUrl
	StateAppOpenURL string = "appOpenUrl"
)

//UIApplicationDelegate ...
type UIApplicationDelegate interface {
	appDidFinishLaunchingWithOptions(options []byte) bool

	appDidRegisterForRemoteNotificationsWithDeviceToken(deviceToken []byte)

	appDidReceiveRemoteNotification(userInfo []byte)

	appDidFailToRegisterForRemoteNotificationsWithError(error []byte)

	appWillResignActive()

	appDidEnterBackground()

	appWillEnterForeground()

	appDidBecomeActive()

	appWillTerminate()

	appOpenUrl(options []byte) bool
}
