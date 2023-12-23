package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
)

// Per vertex phong lighting
func TriangleNormalsPhong(t *graphics.Triangle, viewDirection *basics.Vector3, ambientLightColor *basics.Vector3, specularExponent basics.Scalar, lights []renderLight, specularColor color.Color, ignoreSpecular bool) {
	specularColorAsVector := basics.Vector3FromColor(specularColor)
	for i := 0; i < 3; i++ {
		vertex := &t[i]
		baseColor := vertex.Color
		vertex.Color = ambientTerm(&baseColor, ambientLightColor)
		for _, light := range lights {
			lightVector := light.position.Sub(&vertex.Position)
			lightDistance := lightVector.Length()
			lightVector.ThisNormalize()
			lightFallOff := light.light.FallOff()(lightDistance)
			lightColor := basics.Vector3FromColor(light.light.Color())
			lightColor.ThisMul(lightFallOff)
			vertex.Color.ThisAdd(diffuseTerm(&vertex.Normal, &lightVector, &baseColor, &lightColor))
			finalSpecularColor := specularColorAsVector.Mul(lightFallOff)
			_ = finalSpecularColor

			testLightColor := light.light.Color()
			testLightColorVector := basics.Vector3FromColor(testLightColor)
			if !ignoreSpecular {
				specularTerm := specularTerm(viewDirection, &vertex.Normal, &lightVector, specularExponent, &testLightColorVector, &testLightColorVector)
				vertex.Color.ThisAdd(specularTerm)
			}
		}
		vertex.Color.X = basics.ClampMax(65535, vertex.Color.X)
		vertex.Color.Y = basics.ClampMax(65535, vertex.Color.Y)
		vertex.Color.Z = basics.ClampMax(65535, vertex.Color.Z)
	}
}

// LightFallOff goes from 0 to 1 where 0 is the furthest and 1 is the closest
func PhongLighting(surfaceNormal *basics.Vector3, lightNormal *basics.Vector3, diffuseColor *basics.Vector3, lightColor *basics.Vector3, ambientColor *basics.Vector3, ambientLightColor *basics.Vector3, specularExponent basics.Scalar, specularColor *basics.Vector3, lightFallOff basics.Scalar) basics.Vector3 {

	lightColor.ThisMul(lightFallOff)
	finalColorVector := diffuseTerm(surfaceNormal, lightNormal, diffuseColor, lightColor)

	finalColorVector.ThisAdd(ambientTerm(ambientColor, ambientLightColor))

	//forward := basics.Forward() //hardcoded forward vector for view space
	//specularVector := specularColor.ToVector3()
	//finalColorVector.ThisAdd(specularTerm(&forward, surfaceNormal, lightNormal, specularExponent, &specularVector, &lightColorVector))

	//rescale back to correct range
	//finalColorVector.ThisMul(1 / 65535.0 * 2)
	//finalColorVector.ThisMul(65535)

	finalColorVector.X = basics.ClampMax(65535, finalColorVector.X)
	finalColorVector.Y = basics.ClampMax(65535, finalColorVector.Y)
	finalColorVector.Z = basics.ClampMax(65535, finalColorVector.Z)

	return finalColorVector
}

func diffuseTerm(surfaceNormal *basics.Vector3, lightNormal *basics.Vector3, diffuseColor *basics.Vector3, lightColor *basics.Vector3) basics.Vector3 {
	// (surfaceNormal DOT lightNormal) * ( (diffuse color) per component mul (light color) )
	//(surfaceNormal DOT lightNormal) has to be >= 0

	multiplier := basics.ClampMin(0, surfaceNormal.Dot(lightNormal))

	colorVector := diffuseColor.MulComponents(lightColor)
	colorVector.ThisMul(multiplier)

	//rescale back to correct range
	colorVector.ThisMul(basics.Scalar(1) / basics.Scalar(65535))
	return colorVector
}

func ambientTerm(ambientVector *basics.Vector3, ambientLightVector *basics.Vector3) basics.Vector3 {
	// (ambientColor) per component mul (ambient light color or intensity)

	colorVector := ambientVector.MulComponents(ambientLightVector)
	//rescale back to correct range
	colorVector.ThisMul(basics.Scalar(1) / basics.Scalar(65535))
	return colorVector
}

func specularTerm(viewDirection *basics.Vector3, surfaceNormal *basics.Vector3, lightNormal *basics.Vector3, specularExponent basics.Scalar, specularColor *basics.Vector3, lightColor *basics.Vector3) basics.Vector3 {
	// (surfaceNormal DOT hVersor)^specularExponent * ( (specular color) per component mul (light color) )
	//iLightN := lightNormal.Inverse()
	iViewDir := viewDirection.Inverse()
	hVersor := basics.NLerpVector3(&iViewDir, lightNormal, 0.5)
	multiplier := basics.ClampMin(0, surfaceNormal.Dot(&hVersor))
	if multiplier >= 0.5 {
		//fmt.Println("mult != 0")
	}
	//finalColor := lightColor.MulComponents(specularColor)
	finalColor := *lightColor
	multiplier = basics.Scalar(basics.Pow(multiplier, specularExponent))
	if multiplier >= 0.5 {
		//fmt.Println("mult2 != 0")
	}
	finalColor.ThisMul(multiplier)
	return finalColor
}
