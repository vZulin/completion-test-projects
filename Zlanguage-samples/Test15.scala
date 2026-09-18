final case class User(name: String, active: Boolean)

object ScalaSmoke {
  def greeting(user: User): String = s"Hello, ${user.name} from Scala!"

  def main(args: Array[String]): Unit = {
    val users = List(
      User("IntelliJ IDEA", active = true),
      User("Disabled sample", active = false)
    )

    users.filter(_.active).map(greeting).foreach(println)
  }
}
