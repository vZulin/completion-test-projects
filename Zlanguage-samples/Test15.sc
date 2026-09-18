val languages = List("Scala", "Kotlin", "Java")
val indexedLanguages = languages.zipWithIndex.map {
  case (language, index) => s"${index + 1}: $language"
}

println(indexedLanguages.mkString(System.lineSeparator()))
