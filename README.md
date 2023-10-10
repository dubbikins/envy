# ENVY 👑 💅

Envy is a Go package that provides a flexible way to unmarshal environment variables into struct types using struct tags for validation, type conversion, and more. It offers powerful, customizable tag-based validations like checking for required fields, ensuring values match certain patterns, and verifying they fall within predefined options.

---

## Features
| Feature | Description | Availability |
| :------ |:------: | -----------: |
| **Environment Variable Mapping**  | Easily map environment variables to struct fields using the `env` tag. | Ready ✅ |
| **Defaults** | Define default values for your struct fields using the `default` tag. | Ready ✅ |
| **Options Validation**    | Validate the contents of an environment variable against a set of allowed options using the `options` tag.| Ready ✅ |
| **Regex Matching** | Validate the contents of an environment variable using regular expressions with the `matches` tag. | Ready ✅ |
| **Required Fields** |  Mark certain environment variables as required using the `required` tag.| Ready ✅ |
| **Type Conversion** |  Define your struct with the type you need and envy will take care of converting the environment variable for you.| Ready ✅ |
| **Customizable Tag Parsing Middleware** | Alter the order of tag parsing, add new custom tag parsers, or replace existing ones for enhanced flexibility in the unmarshalling process. | Ready ✅ |
| **Full Feature Test Suite** | Feature tests have been written for Unmarshalling from the env tag | 🧑🏻‍💻 In Progress |
| **Marshalling from struct to environment variables** | This features is on the roadmap, but hadn't been implemented. Open a PR if you'd like to get this feature implemented in a future release! | Not Available 👎 |

![There's a lot of Envious People Hearing That Story](https://media4.giphy.com/media/fwipX8pGKv2VG39w0L/giphy.gif?cid=ecf05e47pz7z69x6o2iqbh8pjoz0priofl97qx7iwtpd68wj&ep=v1_gifs_search&rid=giphy.gif&ct=g) 


## Installation:

```bash
go get -u github.com/dubbikins/envy
```

## Basic Usage

1. **Define Your Struct**:

    Use the `env`,`default`,`options`,`matches`,and `required` tags when defining your struct types.


   ```go
   type Config struct {
      Port     int    `env:"PORT" default:"8080" required:"true"`
      Hostname int    `env:"HOSTNAME" default:"localhost" options:"localhost,my-domain-1.subdomain.com,my-domain-1.subdomain.com"`
      Mode     string `env:"MODE" options:"debug,release,test"`
   }
   ```

1. **Unmarshal Environment Variables**:

   Use the `Unmarshal` function to unmarshal a struct you've already **declared** and **initialized**, just be sure to pass a pointer to the struct, not the struct value. (Note: if you pass an uninitialized value or non-pointer value to Unmarshal, it will return an error)

   ```go
   //HOSTNAME: my-domain-1.subdomain.com
   //Mode: debug

   cfg := &Config{}
   err := envy.Unmarshal(&cfg)
   if err != nil {
      log.Fatal(err)
   }
   ...
   ```

   Or use the *generic* `New` function in conjuction with the `FromEnvironmentAs` function which will return an instance of a type, and unmarshal it with the appropriate environment variables.

   ```go
   if config, err := envy.New(envy.FromEnvironmentAs[Config]);err != nil {
      ...
   }
   ...
   ```

1. **Access the Struct Fields**:

   Then you can utilize the freshly set fields of your struct in the rest of your code! (just be sure to check for errors beforehand 😅)

   ```go
   fmt.Printf("%s:%d [mode=%s]", cfg.Hostname, cfg.Port, cfg.Mode)
   // prints "my-domain-1.subdomain.com:8080 [mode=debug]"
   ```

## Available Type Convertions :boom:

 Envy can unmarshall most primitive go types out of the box. Check out the following table of types to which ones are supported out of the box.

| Type  | Valid Values| Parser|  Provided? |
| ------ | ------ |------ | -----------: |
| `[]byte` `string` |All utf-8 encoded string values `(.*)` | N/A | :white_check_mark: |
| `int` `int8` `int16` `int32` `int64` | positive and negative numbers, optional commas as thousands separators, and no leading zeros: `^([+-]?)(?:[1-9][0-9]{0,2}(?:,[0-9]{3})*\|0)$` |strconv.ParseInt | :white_check_mark: |
| `uint` `uint8``uint16` `uint32` `uint64`| positive numbers, optional commas, as thousands separators and no leading zeros: `^(?:[1-9][0-9]{0,2}(?:,[0-9]{3})*\|0)$` | strconv.ParseUint| :white_check_mark: |
| `float32` `float64` | positive decimal values, commas allowed in the thousands separators, no leading zeros | strcov.ParseFloat | :white_check_mark: |
| `bool` |true,false,yes,no,1,0,on,off,t,f,T,F| strconv.ParseBool (yes,no,on,off pre-parsed to true or false) | :white_check_mark: |
| `embedded structs` `embedded struct pointers` |N/A | envy.Unmarshal | :white_check_mark: |
| `encoding.TextUnmarshaler`  | Whatever you say goes! | Your Custom Implementation | :white_check_mark: |
| `slices` `maps` `time.Time` `time.Duration` `chan`|N/A | require custom | :x: |

### Custom Type unmarshaller 

Didn't find the type you were looking for? :frowning: Not to worry, it's simple to create a custom unmarshaller for any struct type! To provide custom unmarshalling for specific types in your struct, implement the TextUnmarshaler inteface by defining the  `UnmarshalText([]byte) error` method for your type. It's that easy! This interface is a standard go interface from the [encoding package](https://pkg.go.dev/encoding), so you'll find that a lot of custom types already implement this interface.

```go
type CustomType struct {
   // your type definition here
}

func (ct *CustomType) UnmarshalText(data []byte) error {
   // custom parsing logic here
   return nil // or an error
}
```

#### Example: Custom Parsing for Durations :timer_clock:

The below example shows how you can wrap existing types like `Duration` (which is actually just a wrapper for `int64`), from the time package in the go standard library and implement the custom unmarshalling logic to convert the `[]byte` value into a duration value. 

```go

import "time"

// Custom implementation for Parsing Durations from a string 
type Duration struct {
   time.Duration
}

func (d *Duration) UnmarshalText(data []byte) (err error) {
   if len(data) == 0 {
      return nil
   }
   d.Duration, err = time.ParseDuration(string(data))
   return
}
```

Use the custom type in your struct. The below example will convert the environment variable `TIMEOUT` to a time.Duration object as long as the value is in 1,2,5,10, or 25 minutes. If no value is provided, 5 minutes will be used as the defaul duration.

```go
type Config struct {
   Timeout Duration `env:"TIMEOUT" default:"5m" options:"[1m,2m,5m,10m,15m]"`
}
if cfg, err := envy.New(envy.FromEnvironmentAs[Config]);err != nil {
   panic("oh no! something went wrong... 😢")
}
select { 
    case <-time.After(cfg.Timeout.Duration): 
        fmt.Println("Time Out!") 
} 

```

### Reordering the Middleware Stack :repeat:

Envy lets you customize the order of middleware to ensure tag parsers are processed in a sequence of your choice. By utilizing the `Pop` and `Push` methods, you can extract existing middleware and re-insert them or even replace them with custom middleware.

Here's an example:

```go
package main

import (
   "context"
   "fmt"
   "reflect"

   "github.com/your-username/envy"
)

type Config struct {
   // your struct fields with tags here
}

func main() {
   config := &Config{}
   err := envy.Unmarshal(config, func(mw envy.TagMiddleware) {
      // Pop existing middleware
      existingMiddleware1 := mw.Pop()
      existingMiddleware2 := mw.Pop()

      // Push a custom middleware at the start
      mw.Push(
         func(next envy.TagHandler) envy.TagHandler {
            return envy.TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
               fmt.Println("I'm the first custom tag parser")
               return next.UnmarshalField(ctx, field)
            })
         })

      // Re-insert the existing middleware
      mw.Push(existingMiddleware1, existingMiddleware2)

      // Add another custom middleware at the end
      mw.Push(
         func(next envy.TagHandler) envy.TagHandler {
            return envy.TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
               fmt.Println("I'm the last custom tag parser")
               return next.UnmarshalField(ctx, field)
            })
         })
   })

   if err != nil {
      // handle the error
   }
}
```

## Contributions 📣

If you find a bug or think of a new feature, please open an issue or submit a pull request! Go ahead, make the world envious! 😈

![Help Wanted](https://media1.giphy.com/media/25OyfOmwZIARSYKSjL/giphy.gif?cid=ecf05e479gubk5zux0b18o4ualx65ape24kytpnbu31ukyzd&ep=v1_gifs_search&rid=giphy.gif&ct=g)