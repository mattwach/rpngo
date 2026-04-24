// Package key defines key presses.
package key

// ASCII for most keys, except the special ones below
type Key int

const (
	KEY_LEFT      Key = iota + 256
	KEY_RIGHT         // 257
	KEY_UP            // 258
	KEY_DOWN          // 259
	KEY_BACKSPACE     // 260
	KEY_BREAK         // 261
	KEY_DEL           // 262
	KEY_INS           // 263
	KEY_END           // 264
	KEY_HOME          // 265
	KEY_EOF           // 266
	KEY_PAGEUP        // 267
	KEY_PAGEDOWN      // 268
	KEY_F1            // 269
	KEY_F2            // 270
	KEY_F3            // 271
	KEY_F4            // 272
	KEY_F5            // 273
	KEY_F6            // 274
	KEY_F7            // 275
	KEY_F8            // 276
	KEY_F9            // 277
	KEY_F10           // 278
	KEY_F11           // 279
	KEY_F12           // 280
	// Keys pressed with shift held
	KEY_SUP    // 281
	KEY_SDOWN  // 282
	KEY_SLEFT  // 283
	KEY_SRIGHT // 284
	KEY_SHOME  // 285
	KEY_SEND   // 286
	// editing
	KEY_CUT   // 287
	KEY_COPY  // 288
	KEY_PASTE // 289
	KEY_SAVE  // 290
	KEY_QUIT  // 291
	KEY_HELP  // 292
)
