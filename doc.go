/*
Package velty implements subset of the JDK Velocity, using the same syntax as Velocity.
Implemented subset:

variables - i.e. `${foo.Name} $Name`
assignment - i.e. `#set($var1 = 10 + 20 * 10) #set($var2 = ${foo.Name})`
if statements - i.e. `#if(1==1) abc #elsif(2==2) def #else ghi #end`
foreach - i.e. `#foreach($name in ${foo.Names})`
function calls - i.e. `${name.toUpper()}`
template evaluation - i.e. `#evaluate($TEMPLATE)`

*/
package velty
