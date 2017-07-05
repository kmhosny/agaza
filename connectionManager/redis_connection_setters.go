package connectionManager

/*steps for adding a data layer function:
* 1- retrieve a connection from the connections pool
* 2- defer the connection to make sure it's released at the end of the function
* 3- write the data retrieval login
 */
