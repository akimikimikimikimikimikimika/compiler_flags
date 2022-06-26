build:
	@ ( \
		( \
			`flags go macos amd64` -o .~amd64 & \
			`flags go macos arm64` -o .~arm64 & \
			wait \
		) && \
		lipo -create -output bin-macos .~amd64 .~arm64 && \
		rm .~amd64 .~arm64 \
	) & \
	`flags go linux amd64` -o bin-linux-amd64 & \
	`flags go linux arm64` -o bin-linux-arm64 & \
	`flags go windows amd64` -o bin-windows-amd64.exe & \
	wait